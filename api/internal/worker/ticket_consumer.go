package worker

import (
	"classifier/internal/kafka"
	"classifier/internal/tickets"
	"context"
	"log"
	"time"
)

type TicketConsumerWorker struct {
	consumer *kafka.Consumer
	usecase  tickets.TicketClassifierUsecase
}

func NewTicketConsumerWorker(consumer *kafka.Consumer, usecase tickets.TicketClassifierUsecase) *TicketConsumerWorker {
	return &TicketConsumerWorker{
		consumer: consumer,
		usecase:  usecase,
	}
}

func (w *TicketConsumerWorker) Run(ctx context.Context) {
	for {
		key, value, err := w.consumer.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("Received message from Kafka, key=%s", string(key))

		if err := w.usecase.ProcessClassifiedTicket(ctx, value); err != nil {
			log.Printf("failed to process classified ticket: %v", err)
		}
	}
}
