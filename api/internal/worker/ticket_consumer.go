package worker

import (
	"classifier/internal/kafka"
	"classifier/internal/tickets"
	"context"
	"errors"
	"io"
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
		select {
		case <-ctx.Done():
			log.Println("TicketConsumerWorker stop")
			return
		default:
			key, value, err := w.consumer.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
					log.Println("Consumer stopped gracefully")
					return
				}
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
}
