package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
	"visit-service/internal/models"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

var unknownLocation = &models.IpLocationResponse{
	Status:     "UNKNOWN",
	Country:    "UNKNOWN",
	RegionName: "UNKNOWN",
	City:       "UNKNOWN",
}

func GetLocation(ip string) (*models.IpLocationResponse, error) {
	url := os.Getenv("MERIDIAN_URL") + "/location/" + ip

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("Error creating request to ip-api", "ip", ip, "error", err)
		return unknownLocation, nil
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			slog.Error("Timeout getting location from ip-api", "ip", ip, "timeout", "30s")
		} else {
			slog.Error("Error getting location from ip-api", "ip", ip, "error", err)
		}
		return unknownLocation, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Error getting location from ip-api", "ip", ip, "status", resp.StatusCode)
		return unknownLocation, nil
	}

	var location models.IpLocationResponse

	err = json.NewDecoder(resp.Body).Decode(&location)
	if err != nil {
		slog.Error("Error decoding ip-api response", "error", err)
		return unknownLocation, nil
	}
	return &location, nil
}
