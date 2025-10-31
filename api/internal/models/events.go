package models

import "time"

type ProcessedEvents struct {
	EventID          string
	TicketID    int64
	ProcessedAt time.Time
}
