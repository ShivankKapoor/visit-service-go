package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"visit-service/internal/models"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func GetLocation(ip string) (*models.IpLocationResponse, error) {
	url := "http://ip-api.com/json/" + ip + "?fields=status,country,regionName,city"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("Error creating request to ip-api", "ip", ip, "error", err)
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			slog.Error("Timeout getting location from ip-api", "ip", ip, "timeout", "30s")
		} else {
			slog.Error("Error getting location from ip-api", "ip", ip, "error", err)
		}
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Status Code %d\n", resp.StatusCode)
		return nil, fmt.Errorf("ip-api returned status %d", resp.StatusCode)
	}

	var location models.IpLocationResponse

	err = json.NewDecoder(resp.Body).Decode(&location)
	if err != nil {
		slog.Error("Error decoding ip-api response", "error", err)
		return nil, err
	}
	return &location, nil
}
