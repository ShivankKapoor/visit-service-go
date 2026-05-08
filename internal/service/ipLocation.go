package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"visit-service/internal/models"
)

func GetLocation(ip string) (*models.IpLocationResponse, error) {
	url := "http://ip-api.com/json/" + ip + "?fields=status,country,regionName,city"

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Error getting location from ip-api", "ip", ip, "error", err)
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
