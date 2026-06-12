package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"visit-service/internal/dto"
	"visit-service/internal/service"
)

type TrackHandler struct{}

func NewTrackHandler() *TrackHandler {
	return &TrackHandler{}
}

func (h *TrackHandler) Track(w http.ResponseWriter, r *http.Request) {
	slog.Info("Track request incoming")
	var req dto.TrackRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.DeviceInfo == "" || req.PageVisited == "" {
		http.Error(w, "deviceInfo and pageVisited are required", http.StatusBadRequest)
		return
	}
	service.Track(req, r)
}
