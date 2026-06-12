package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"visit-service/internal/dto"
)

type MainHandler struct{}

func NewMainHandler() *MainHandler {
	return &MainHandler{}
}

func (h *MainHandler) Home(w http.ResponseWriter, r *http.Request) {
	slog.Info("Home endpoint called")

	resp := &dto.HomeReponseDTO{
		Name:     "Visit Service",
		Platform: "Go",
		Status:   "Up",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
