package service

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
	"visit-service/internal/dto"
)

var locationClient = &http.Client{
	Timeout: 10 * time.Second,
}

var unknownLocation = dto.MeridianResponseDTO{
	Country:     "UNKNOWN",
	CountryCode: "UNKNOWN",
	City:        "UNKNOWN",
	RegionName:  "UNKNOWN",
}

func GetLocation(ip string) dto.MeridianResponseDTO {
	meridianURL := os.Getenv("MERIDIAN_URL")
	if meridianURL == "" {
		slog.Error("Meridian URL not set")
		return unknownLocation
	}

	slog.Info("Getting location", "ip", ip)
	requestURL := meridianURL + "/location/" + ip

	resp, err := locationClient.Get(requestURL)
	if err != nil {
		slog.Error("Failed to get location via meridian", "err", err)
		return unknownLocation
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Error("Meridian returned non-success status", "status", resp.StatusCode)
		return unknownLocation
	}

	var location dto.MeridianResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		slog.Error("Failed to decode location response", "err", err)
		return unknownLocation
	}

	return location
}
