package handler

import (
	"net/http"
	"os"
	"visit-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CronHandler struct {
	db *pgxpool.Pool
}

func NewCronHandler(db *pgxpool.Pool) *CronHandler {
	return &CronHandler{db: db}
}

func (h *CronHandler) RunDailySummary(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("PROD") == "true" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	service.RunDailySummary(h.db)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("daily summary ran"))
}
