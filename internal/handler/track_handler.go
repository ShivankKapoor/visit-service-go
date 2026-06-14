package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"visit-service/internal/dto"
	"visit-service/internal/service"
)

type TrackHandler struct {
	trackService *service.TrackService
}

func NewTrackHandler(trackService *service.TrackService) *TrackHandler {
	return &TrackHandler{trackService: trackService}
}

func (h *TrackHandler) Track(w http.ResponseWriter, r *http.Request) {

	slog.Info("Track request incoming")
	var req dto.TrackRequest

	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.DeviceInfo == "" || req.PageVisited == "" {
		http.Error(w, "deviceInfo and pageVisited are required", http.StatusBadRequest)
		return
	}

	ip := service.GetClientIP(r)
	userAgent := r.Header.Get("User-Agent")

	go h.trackService.Track(req, ip, userAgent)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("OK")
}
