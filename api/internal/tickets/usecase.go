package tickets

import (
	"classifier/internal/models"
	"classifier/internal/outbox"
	"context"
	"encoding/json"
)

type TicketUsecase interface {
	GetTickets(ctx context.Context) ([]models.Ticket, error)
	CreateTicket(ctx context.Context, title string) error
}

type ticketUsecase struct {
	repo       TicketRepository
	outboxRepo outbox.OutboxRepository
}

func NewTicketUsecase(r TicketRepository, or outbox.OutboxRepository) TicketUsecase {
	return &ticketUsecase{repo: r, outboxRepo: or}
}

func (u *ticketUsecase) GetTickets(ctx context.Context) ([]models.Ticket, error) {
	return u.repo.GetAll(ctx)
}

func (u *ticketUsecase) CreateTicket(ctx context.Context, text string) error {
	ticket := &models.Ticket{
		Text:     text,
		Category: "pending",
	}

	event := map[string]any{
		"text": text,
	}
	payload, _ := json.Marshal(event)

	return u.repo.CreateWithOutbox(ctx, ticket, "ticket_created", payload)
}
