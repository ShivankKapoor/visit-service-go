package models

type HealthResponse struct {
	Name      string `json:"name"`
	Platform  string `json:"platform"`
	APIStatus string `json:"api-status"`
	DBStatus  string `json:"db-status"`
}
