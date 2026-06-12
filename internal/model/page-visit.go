package model

type PageVisit struct {
	ID          string
	IPAddress   string
	PageVisited string
	DeviceInfo  *string
	UserAgent   *string
	Timestamp   string
}
