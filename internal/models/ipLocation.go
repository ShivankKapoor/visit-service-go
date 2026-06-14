package models

type IpLocationResponse struct {
	Status     string `json:"status"`
	Country    string `json:"country"`
	RegionName string `json:"regionName"`
	City       string `json:"city"`
}
