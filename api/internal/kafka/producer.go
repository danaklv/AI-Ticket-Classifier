package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokerURL, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: w}
}

func (p *Producer) SendMessage(ctx context.Context, key string, value []byte) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
	if err != nil {
		log.Printf("failed to write message: %v", err)
	}
	return err
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
