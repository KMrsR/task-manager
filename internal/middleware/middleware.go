package middleware

import (
	"net/http"
	"strings"
)

// проверка для POST/PUT что в теле запроса есть JSON
func RequireJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			handlers.responseWithError(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		}

		next.ServeHTTP(w, r)
	})
}
