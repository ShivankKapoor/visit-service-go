package models

import "time"

type DailyVisitStats struct {
	SummaryDate time.Time
	TotalVisits int64
}
