package models

type HealthCheckResponse struct {
	Name     string `json:"name"`
	Platform string `json:"platform"`
	Status   string `json:"status"`
}
