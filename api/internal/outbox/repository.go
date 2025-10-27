package outbox

import (
	"classifier/internal/models"
	"context"

	"gorm.io/gorm"
)

type OutboxRepository interface {
	AddEvent(ctx context.Context, eventType string, payload []byte) error
	GetUnsent(ctx context.Context, limit int) ([]models.Outbox, error)
	MarkAsSent(ctx context.Context, id uint) error
}

type outboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) OutboxRepository {
	return &outboxRepository{db: db}
}

func (r *outboxRepository) AddEvent(ctx context.Context, eventType string, payload []byte) error {
	outbox := &models.Outbox{
		EventType: eventType,
		Payload:   payload,
	}
	return r.db.WithContext(ctx).Create(outbox).Error
}

func (r *outboxRepository) GetUnsent(ctx context.Context, limit int) ([]models.Outbox, error) {
	var events []models.Outbox
	err := r.db.WithContext(ctx).
		Where("sent = ?", false).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *outboxRepository) MarkAsSent(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.Outbox{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"sent": true,
		}).Error
}
