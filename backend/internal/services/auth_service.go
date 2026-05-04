package services

import (
	"os"
	"strings"
	"time"

	"ai-inference-gateway/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultRegistrationBalance = 100
	tokenTTL                   = 24 * time.Hour
)

type AuthService struct {
	userRepo UserRepository
}

type AuthClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
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
		return nil, "", ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if isRepoNotFoundError(err, "user not found:") {
			return nil, "", ErrInvalidCredentials
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
	return token.SignedString(jwtSecret())
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
		return jwtSecret(), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.UserID == "" || claims.Email == "" || claims.Role == "" {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func validateAuthInput(username, email, password string) error {
	if username == "" {
		return ErrInvalidAuthInput
	}
	if email == "" || !strings.Contains(email, "@") {
		return ErrInvalidAuthInput
	}
	if len(password) < 6 {
		return ErrInvalidAuthInput
	}

	return nil
}

func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Development fallback for local labs only. Set JWT_SECRET outside development.
		secret = "dev-secret"
	}

	return []byte(secret)
}
