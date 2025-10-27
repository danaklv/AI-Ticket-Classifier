package tickets

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type TicketClassifierUsecase interface {
	ProcessClassifiedTicket(ctx context.Context, data []byte) error
}

type ticketClassifierUsecase struct {
	repo TicketRepository
}

func NewTicketClassifierUsecase(repo TicketRepository) TicketClassifierUsecase {
	return &ticketClassifierUsecase{repo: repo}
}

type classifiedTicketEvent struct {
	ID                int64  `json:"ticket_id"`
	PredictedCategory string `json:"predicted_category"`
}

func (u *ticketClassifierUsecase) ProcessClassifiedTicket(ctx context.Context, data []byte) error {
	var event classifiedTicketEvent
	fmt.Println(string(data))
	if err := json.Unmarshal(data, &event); err != nil {
		log.Printf("failed to unmarshal classified ticket: %v", err)
		return err
	}

	log.Printf("Updating ticket %d with category: %s", event.ID, event.PredictedCategory)
	return u.repo.UpdateCategoryByID(ctx, event.ID, event.PredictedCategory)
}
