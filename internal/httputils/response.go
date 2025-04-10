package httputils

import (
	"encoding/json"
	"net/http"
)

// функция дл стандартных ответов
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// функция для ответа с ошибкой
func responseWithError(w http.ResponseWriter, msg string, code int) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}
