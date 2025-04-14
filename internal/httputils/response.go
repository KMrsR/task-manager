package httputils

import (
	"encoding/json"
	"net/http"
)

// функция для стандартных ответов
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// функция для ответа с ошибкой
func ResponseWithError(w http.ResponseWriter, msg string, code int) {
	RespondWithJSON(w, code, map[string]string{"error": msg})
}
