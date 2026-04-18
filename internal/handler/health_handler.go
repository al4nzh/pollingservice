package handler

import (
	"encoding/json"
	"net/http"

	"github.com/al4nzh/pollingservice.git/internal/service"
)

type HealthHandler struct {
	healthService *service.HealthService
}

func NewHealthHandler(healthService *service.HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

func (h *HealthHandler) GetHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.healthService.GetHealth())
}

func (h *HealthHandler) GetReady(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.healthService.GetReadiness())
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}
