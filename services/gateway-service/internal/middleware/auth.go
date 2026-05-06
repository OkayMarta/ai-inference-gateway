package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authUserContextKey struct{}

type AuthUser struct {
	ID    string
	Email string
	Role  string
}

type authClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func JWT(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenValue, ok := bearerToken(r.Header.Get("Authorization"))
			if !ok {
				respondError(w, r, http.StatusUnauthorized, "unauthorized")
				return
			}

			claims := &authClaims{}
			token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenSignatureInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid || claims.UserID == "" || claims.Email == "" || claims.Role == "" {
				respondError(w, r, http.StatusUnauthorized, "unauthorized")
				return
			}

			user := &AuthUser{ID: claims.UserID, Email: claims.Email, Role: claims.Role}
			ctx := context.WithValue(r.Context(), authUserContextKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func respondError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(struct {
		Timestamp string `json:"timestamp"`
		Status    int    `json:"status"`
		Message   string `json:"message"`
		Path      string `json:"path"`
	}{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Status:    status,
		Message:   message,
		Path:      r.URL.Path,
	})
}

func CurrentUser(r *http.Request) (*AuthUser, bool) {
	user, ok := r.Context().Value(authUserContextKey{}).(*AuthUser)
	return user, ok
}

func bearerToken(authorization string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorization, prefix))
	return token, token != ""
}
