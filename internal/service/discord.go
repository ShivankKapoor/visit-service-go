package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
	"visit-service/internal/models"
)

var discordClient = &http.Client{
	Timeout: 30 * time.Second,
}

func send(req models.DiscordRequest) error {
	dsn := os.Getenv("DISCORD_WEBHOOK_URL")
	if dsn == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL not set")
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error marshaling Discord message", "error", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", dsn, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Error creating Discord request", "error", err)
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := discordClient.Do(httpReq)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			slog.Error("Timeout sending Discord message", "timeout", "30s")
		} else {
			slog.Error("Error sending Discord message", "error", err)
		}
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		slog.Warn("Discord API returned non-success status", "status", resp.StatusCode)
	}

	return nil
}

func getEmoji() string {
	customEmoji := os.Getenv("CUSTOM_EMOJI_ID")
	emoji := "🐹"
	if customEmoji != "" {
		emoji = customEmoji
	}
	return emoji
}

func SendVisitMessage(visit models.PageVisit, location string) {
	device := "unknown"
	if visit.DeviceInfo != nil {
		device = *visit.DeviceInfo
	}

	content := fmt.Sprintf("🌎 Visitor\nIP: %s\nLocation: %s\nPage: %s\nDevice: %s\n%s\n", visit.IPAddress, location, visit.PageVisited, device, getEmoji())

	body := models.DiscordRequest{
		Content: content,
	}

	send(body)
}

func SendDailyVisitsMessage(noOfVisits int) {
	content := fmt.Sprintf("✅ Cron Report \nNumber of Visits: %d\n%s\n", noOfVisits, getEmoji())

	body := models.DiscordRequest{
		Content: content,
	}

	send(body)
}
