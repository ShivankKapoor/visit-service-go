package handler

import (
	"log/slog"
	"net/http"
)

type TrackHandler struct{}

func NewTrackHandler() *TrackHandler {
	return &TrackHandler{}
}

func (h *TrackHandler) Track(w http.ResponseWriter, r *http.Request) {
	slog.Info("Track request incoming")

}
