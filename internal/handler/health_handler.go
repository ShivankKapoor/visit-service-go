package handler

import (
	"encoding/json"
	"net/http"
	"visit-service/internal/dto"
	"visit-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	dbStatus := "Healthy"
	statusCode := http.StatusOK

	if !service.GetDBHealth(h.db) {
		dbStatus = "Unhealthy"
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.HealthResponseDTO{
		Name:      "Visit Service",
		Platform:  "GO",
		APIStatus: "UP",
		DBStatus:  dbStatus,
	})
}
