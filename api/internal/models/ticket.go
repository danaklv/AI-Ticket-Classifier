package models

type Ticket struct {
	ID       int64  `json:"id"`
	Text     string `json:"text"`
	Category string `json:"category"`
}

type ClassifiedTicket struct {
	ID                any    `json:"id"`
	Subject           string `json:"subject"`
	Body              string `json:"body"`
	PredictedCategory string `json:"predicted_category"`
	ProcessedAt       string `json:"processed_at"`
}
