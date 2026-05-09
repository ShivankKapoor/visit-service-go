package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"visit-service/internal/models"
)

func send(req models.DiscordRequest) error {
	dsn := os.Getenv("DISCORD_WEBHOOK_URL")
	if dsn == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_URL not set")
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(dsn, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
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

	content := fmt.Sprintf("🌎 Visitor\nIP: %s\nLocation: %s\nPage: %s\nDevice: %s\n%s", visit.IPAddress, location, visit.PageVisited, device, getEmoji())

	body := models.DiscordRequest{
		Content: content,
	}

	send(body)
}
