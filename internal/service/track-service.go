package service

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
	"visit-service/internal/dto"
	"visit-service/internal/model"
	"visit-service/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackService struct {
	db *pgxpool.Pool
}

func NewTrackService(db *pgxpool.Pool) *TrackService {
	return &TrackService{db: db}
}

func (s *TrackService) Track(dto dto.TrackRequest, r *http.Request) {
	ip := GetClientIP(r)
	userAgent := r.Header.Get("User-Agent")

	deviceInfo := dto.DeviceInfo
	if idx := strings.Index(deviceInfo, ","); idx != -1 {
		deviceInfo = deviceInfo[:idx]
	}

	visit := model.PageVisit{
		ID:          uuid.New().String(),
		IPAddress:   ip,
		PageVisited: dto.PageVisited,
		DeviceInfo:  &deviceInfo,
		UserAgent:   &userAgent,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	slog.Info("Track request parsed", "id", visit.ID, "ip", visit.IPAddress, "device", *visit.DeviceInfo, "userAgent", *visit.UserAgent, "timestamp", visit.Timestamp)

	visitRepository := repository.NewPageVisitRepository(s.db)
	if err := visitRepository.InsertPageVisit(r.Context(), visit); err != nil {
		slog.Error("Failed to insert page visit", "error", err)
	}

	location := GetLocation(visit.IPAddress)
	locationString := (location.City + ", " + location.RegionName + ", " + location.Country)

	SendVisitMessage(visit, locationString)

}
