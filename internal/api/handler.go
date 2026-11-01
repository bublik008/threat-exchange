package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bublik008/threat-exchange/internal/api/gen"
	"github.com/bublik008/threat-exchange/internal/repository"
	"github.com/bublik008/threat-exchange/internal/service"
)

// Server реализует сгенерированный интерфейс gen.ServerInterface.
type Server struct {
	svc *service.ThreatService
}

// NewServer создает новый HTTP-сервер.
func NewServer(svc *service.ThreatService) *Server {
	return &Server{svc: svc}
}

// CreateThreat обрабатывает POST /info
func (s *Server) CreateThreat(w http.ResponseWriter, r *http.Request) {
	var req gen.ThreatCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	resp, err := s.svc.CreateThreat(r.Context(), &req)
	if err != nil {
		s.handleError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, resp)
}

// GetInfoByMd5 обрабатывает GET /info/{md5}
func (s *Server) GetInfoByMd5(w http.ResponseWriter, r *http.Request, md5 string) {
	resp, err := s.svc.GetThreat(r.Context(), md5)
	if err != nil {
		s.handleError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// UpdateThreat обрабатывает PUT /info/{md5}
func (s *Server) UpdateThreat(w http.ResponseWriter, r *http.Request, md5 string) {
	var req gen.ThreatUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	resp, err := s.svc.UpdateThreat(r.Context(), md5, &req)
	if err != nil {
		s.handleError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

// handleError — централизованный маппинг доменных ошибок в HTTP-статусы.
func (s *Server) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidMD5), errors.Is(err, service.ErrInvalidInput):
		respondError(w, http.StatusBadRequest, "Validation error", err.Error())
	case errors.Is(err, repository.ErrNotFound):
		respondError(w, http.StatusNotFound, "Threat not found", err.Error())
	case errors.Is(err, repository.ErrConflict):
		respondError(w, http.StatusConflict, "Threat already exists", err.Error())
	default:
		// Любые другие ошибки (например, упала БД) отдаем как 500.
		respondError(w, http.StatusInternalServerError, "Internal server error", "")
	}
}