package service

import (
	"log/slog"
	"net/http"
	"visit-service/internal/dto"
)

func Track(dto dto.TrackRequest, r *http.Request) {
	var ip = GetClientIP(r)
	slog.Info("Track request parsed", "deviceInfo", dto.DeviceInfo, "pageVisited", dto.PageVisited, "Ip", ip)
}
