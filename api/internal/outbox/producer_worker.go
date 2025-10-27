package outbox

import (
	"classifier/internal/kafka"
	"context"
	"log"
	"time"
)

type OutboxWorker struct {
	repo     OutboxRepository
	producer *kafka.Producer
}

func NewOutboxWorker(repo OutboxRepository, producer *kafka.Producer) *OutboxWorker {
	return &OutboxWorker{repo: repo, producer: producer}
}

func (w *OutboxWorker) Run(ctx context.Context) {
	for {
	
		events, err := w.repo.GetUnsent(ctx, 10)
		if err != nil {
			log.Println("Ошибка получения событий:", err)
			time.Sleep(5 * time.Second)
			continue
		}


		for _, e := range events {
			err := w.producer.SendMessage(ctx, e.EventType, e.Payload)
			if err != nil {
				log.Printf("Ошибка отправки события %d: %v", e.ID, err)
				continue
			}

		
			if err := w.repo.MarkAsSent(ctx, e.ID); err != nil {
				log.Printf("Ошибка обновления статуса %d: %v", e.ID, err)
			}
		}

	
		time.Sleep(3 * time.Second)
	}
}
