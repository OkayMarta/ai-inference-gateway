package handlers

import (
	"crypto/subtle"
	"net/http"
)

const internalServiceTokenHeader = "X-Internal-Service-Token"

func InternalServiceTokenMiddleware(expectedToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(internalServiceTokenHeader)
			if token == "" || subtle.ConstantTimeCompare([]byte(token), []byte(expectedToken)) != 1 {
				respondError(w, r, http.StatusUnauthorized, "unauthorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
