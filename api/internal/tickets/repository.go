package tickets

import (
	"classifier/internal/models"
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

type TicketRepository interface {
	GetAll(ctx context.Context) ([]models.Ticket, error)
	Create(ctx context.Context, t *models.Ticket) error
	CreateWithOutbox(ctx context.Context, t *models.Ticket, eventType string, payload []byte) error
	UpdateCategoryByID(ctx context.Context, id int64, category string) error
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

func (r *ticketRepository) UpdateCategoryByID(ctx context.Context, id int64, category string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(&models.Ticket{}).
			Where("id = ?", id).
			Update("category", category).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *ticketRepository) CreateWithOutbox(ctx context.Context, t *models.Ticket, eventType string, payload []byte) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(t).Error; err != nil {
			return err
		}
		event := map[string]any{
			"id":   t.ID,
			"text": t.Text,
		}
		payload, _ := json.Marshal(event)

		outbox := &models.Outbox{
			EventType: eventType,
			Payload:   payload,
		}
		if err := tx.Create(outbox).Error; err != nil {
			return err
		}

		return nil
	})
}
