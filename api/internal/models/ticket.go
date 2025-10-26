package models


type Ticket struct {
	ID    int64  `json:"id"`
	Text string `json:"text"`
	Category string `json:"category"`
}
