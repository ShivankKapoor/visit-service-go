package models

type TrackRequest struct {
	DeviceInfo  string `json:"deviceInfo"`
	PageVisited string `json:"pageVisited"`
}

type PageVisit struct {
	ID          string
	IPAddress   string
	PageVisited string
	DeviceInfo  *string
	UserAgent   *string
	Timestamp   string
}
