package service

import (
	"context"
	"log/slog"
	"visit-service/internal/models"
	"visit-service/internal/repositories"
)

func TrackAsync(visit models.PageVisit) {
	err := repositories.InsertPageVisit(context.Background(), visit)
	if err != nil {
		slog.Error("Failed to insert page visit", "error", err)
		return
	}

	locationReq, err := GetLocation(visit.IPAddress)
	location := "unknown"
	if err == nil {
		location = locationReq.City + ", " + locationReq.RegionName + ", " + locationReq.Country
	}

	slog.Info("Visit recorded", "ip", visit.IPAddress, "page", visit.PageVisited, "device", visit.DeviceInfo, "location", location)
}
