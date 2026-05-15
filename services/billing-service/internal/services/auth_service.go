package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
	"unicode"

	"billing-service/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultRegistrationBalance = 100
	tokenTTL                   = 24 * time.Hour
)

type AuthService struct {
	userRepo    UserRepository
	resetRepo   PasswordResetRepository
	emailSender PasswordResetEmailSender
	jwtSecret   []byte
}

type AuthClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo UserRepository, resetRepo PasswordResetRepository, emailSender PasswordResetEmailSender, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		resetRepo:   resetRepo,
		emailSender: emailSender,
		jwtSecret:   []byte(jwtSecret),
	}
}

func (s *AuthService) Register(username, email, password string) (*models.User, string, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if err := validateAuthInput(username, email, password); err != nil {
		return nil, "", err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		ID:           generateID(),
		Username:     username,
		Email:        email,
		PasswordHash: string(passwordHash),
		TokenBalance: defaultRegistrationBalance,
		Role:         models.RoleUser,
	}

	if err := s.userRepo.Create(user); err != nil {
		if isDuplicateEmailError(err) {
			return nil, "", ErrEmailAlreadyExists
		}
		return nil, "", err
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	email = strings.TrimSpace(email)
	if email == "" || !strings.Contains(email, "@") || password == "" {
		return nil, "", ErrInvalidLoginInput
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if isRepoNotFoundError(err, "user not found:") {
			return nil, "", ErrAccountNotFound
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) RequestPasswordReset(email string) error {
	email = strings.TrimSpace(email)
	if email == "" || !strings.Contains(email, "@") {
		return ErrInvalidEmail
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if isRepoNotFoundError(err, "user not found:") {
			return nil
		}
		return err
	}

	rawToken, err := generateResetToken()
	if err != nil {
		return err
	}

	token := &models.PasswordResetToken{
		ID:        generateID(),
		UserID:    user.ID,
		TokenHash: hashResetToken(rawToken),
		ExpiresAt: time.Now().UTC().Add(time.Duration(passwordResetTTLMinutes()) * time.Minute),
	}

	if err := s.resetRepo.Create(token); err != nil {
		return err
	}

	resetLink, err := buildResetLink(rawToken)
	if err != nil {
		return err
	}

	if err := s.emailSender.SendPasswordResetEmail(user.Email, resetLink); err != nil {
		log.Printf("failed to send password reset email for user %s: %v", user.ID, err)
		return nil
	}

	return nil
}

func (s *AuthService) ResetPassword(tokenValue, newPassword string) error {
	tokenValue = strings.TrimSpace(tokenValue)
	if tokenValue == "" {
		return ErrInvalidPasswordResetInput
	}
	if err := validatePassword(newPassword); err != nil {
		return err
	}

	resetToken, err := s.resetRepo.GetValidByTokenHash(hashResetToken(tokenValue))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "password reset token not found") {
			return ErrInvalidPasswordResetToken
		}
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePasswordHash(resetToken.UserID, string(passwordHash)); err != nil {
		return err
	}

	if err := s.resetRepo.MarkUsed(resetToken.ID); err != nil {
		return fmt.Errorf("password updated but reset token could not be marked used: %w", err)
	}

	return nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	if user == nil || user.ID == "" || user.Email == "" || user.Role == "" {
		return "", ErrInvalidAuthInput
	}

	now := time.Now().UTC()
	claims := AuthClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenValue string) (*AuthClaims, error) {
	tokenValue = strings.TrimSpace(tokenValue)
	if tokenValue == "" {
		return nil, ErrInvalidToken
	}

	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.UserID == "" || claims.Email == "" || claims.Role == "" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func validateAuthInput(username, email, password string) error {
	if username == "" {
		return ErrUsernameRequired
	}
	if email == "" || !strings.Contains(email, "@") {
		return ErrInvalidEmail
	}
	if err := validatePassword(password); err != nil {
		return err
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}

	hasLetter := false
	hasDigit := false
	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return ErrWeakPassword
	}

	return nil
}

func generateResetToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("generate password reset token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(tokenBytes), nil
}

func hashResetToken(tokenValue string) string {
	hash := sha256.Sum256([]byte(tokenValue))
	return hex.EncodeToString(hash[:])
}

func buildResetLink(tokenValue string) (string, error) {
	frontendURL := envOrDefault("FRONTEND_URL", "http://localhost:5173")
	parsed, err := url.Parse(frontendURL)
	if err != nil {
		return "", fmt.Errorf("parse frontend url: %w", err)
	}

	parsed.Path = "/reset-password"
	query := parsed.Query()
	query.Set("token", tokenValue)
	parsed.RawQuery = query.Encode()

	return parsed.String(), nil
}

func isDuplicateEmailError(err error) bool {
	if err == nil {
		return false
	}

	errText := strings.ToLower(err.Error())
	return strings.Contains(errText, "users_email_unique") ||
		strings.Contains(errText, "users_email_key") ||
		strings.Contains(errText, "duplicate key") ||
		strings.Contains(errText, "email already exists")
}
