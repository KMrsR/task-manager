package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KMrsR/task-manager/internal/httputils"
)

// проверка для POST/PUT что в теле запроса есть JSON
func RequireJSON(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			httputils.ResponseWithError(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// логирование времени выполнения и статускодов ответов
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		//обертка для responsewriter
		lrw := &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		log.Printf("%s %s -> %d (%v)",
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			time.Since(start))

	})
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *logResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
