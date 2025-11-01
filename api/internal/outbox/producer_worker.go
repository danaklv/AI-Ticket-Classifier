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

	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("OutboxWorker stop")
			return
		case <-t.C:
			events, err := w.repo.GetUnsent(ctx, 10)
			if err != nil {
				log.Println("get unsent:", err)

				continue
			}
			for _, e := range events {
				err := w.producer.SendMessage(ctx, e.EventType, e.Payload)
				if err != nil {
					log.Printf("send message %d: %v", e.ID, err)
					continue
				}

				if err := w.repo.MarkAsSent(ctx, e.ID); err != nil {
					log.Printf("mark %d: %v", e.ID, err)
				}
			}
		}
	}
}
