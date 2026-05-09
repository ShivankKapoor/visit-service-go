package service

import (
	"context"
	"log/slog"
	"time"
	"visit-service/internal/models"
	"visit-service/internal/repositories"
)

func TrackAsync(visit models.PageVisit) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := repositories.InsertPageVisit(ctx, visit)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			slog.Error("Timeout inserting page visit", "ip", visit.IPAddress, "timeout", "30s")
		} else {
			slog.Error("Failed to insert page visit", "ip", visit.IPAddress, "error", err)
		}
		return
	}

	locationReq, err := GetLocation(visit.IPAddress)
	location := "unknown"
	if err == nil {
		location = locationReq.City + ", " + locationReq.RegionName + ", " + locationReq.Country
	}
	SendVisitMessage(visit, location)
	slog.Info("Visit recorded", "ip", visit.IPAddress, "page", visit.PageVisited, "device", visit.DeviceInfo, "location", location)
}
