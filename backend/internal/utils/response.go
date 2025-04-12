// File: /internal/utils/response.go

package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse estrutura um erro padronizado para API
type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

// SendError retorna um erro JSON padronizado
func SendError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
		Error:   http.StatusText(status),
		Code:    status,
	})
}

// SendSuccess retorna uma resposta JSON padronizada
func SendSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// SendJSON retorna uma resposta JSON com o status especificado
func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
