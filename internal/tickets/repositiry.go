package tickets

import (
	"classifier/internal/models"
	"context"

	"gorm.io/gorm"
)

type TicketRepository interface {
	GetAll(ctx context.Context) ([]models.Ticket, error)
	Create(ctx context.Context, t *models.Ticket) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) GetAll(ctx context.Context) ([]models.Ticket, error) {

	var tickets []models.Ticket
	rows := r.db.Find(&tickets)

	if rows.Error != nil {
		return nil, rows.Error
	}

	return tickets, nil
}

func (r *ticketRepository) Create(ctx context.Context, t *models.Ticket) error {
	res := r.db.Create(t)

	return res.Error
}
