package tickets

import (
	"classifier/internal/models"
	"context"
)

type TicketUsecase interface {
	GetTickets(ctx context.Context) ([]models.Ticket, error)
	CreateTicket(ctx context.Context, title string) error
}

type ticketUsecase struct {
	repo TicketRepository
}

func NewTicketUsecase(r TicketRepository) TicketUsecase {
	return &ticketUsecase{repo: r}
}

func (u *ticketUsecase) GetTickets(ctx context.Context) ([]models.Ticket, error) {
	return u.repo.GetAll(ctx)
}

func (u *ticketUsecase) CreateTicket(ctx context.Context, text string) error {
	return u.repo.Create(ctx, &models.Ticket{Text: text})
}
