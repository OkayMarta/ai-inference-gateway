package handlers

import (
	"encoding/json"
	"net/http"
)

// APIError — це єдиний формат для всіх помилок нашого API. Завдяки цьому фронтенд завжди знає, яку структуру очікувати при збоях
type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// respondJSON — помічник для успішних відповідей. Він автоматично ставить заголовок Content-Type та перетворює дані у JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError — помічник для відправки помилок
func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, APIError{
		Error:   http.StatusText(status),
		Message: msg,
		Status:  status,
	})
}

// RecoveryMiddleware: Якщо десь у коді станеться "panic" (наприклад, звернення до nil вказівника), цей middleware перехопить паніку і не дасть серверу впасти повністю. Замість крашу він поверне клієнту 500 Internal Server Error
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				respondError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}