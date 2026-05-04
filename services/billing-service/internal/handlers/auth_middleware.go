package handlers

import (
	"context"
	"net/http"
	"strings"

	"billing-service/internal/services"
)

type authUserContextKey struct{}

type AuthUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func AuthMiddleware(auth *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenValue, ok := bearerToken(r.Header.Get("Authorization"))
			if !ok {
				respondError(w, r, http.StatusUnauthorized, services.ErrUnauthorized.Error())
				return
			}

			claims, err := auth.ValidateToken(tokenValue)
			if err != nil {
				respondError(w, r, http.StatusUnauthorized, services.ErrUnauthorized.Error())
				return
			}

			user := &AuthUser{
				ID:    claims.UserID,
				Email: claims.Email,
				Role:  claims.Role,
			}

			ctx := context.WithValue(r.Context(), authUserContextKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
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
