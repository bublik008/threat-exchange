// Package api реализует HTTP-транспорт для сервиса.
package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bublik008/threat-exchange/internal/api/gen"
)

// respondJSON пишет данные в формате JSON с указанным HTTP-статусом.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			slog.Error("failed to encode response", "error", err)
		}
	}
}

// respondError формирует стандартный ответ с ошибкой согласно OpenAPI-спеке.
func respondError(w http.ResponseWriter, status int, message string, details string) {
	respondJSON(w, status, gen.ErrorResponse{
		Message: &message,
		Details: &details,
	})
}