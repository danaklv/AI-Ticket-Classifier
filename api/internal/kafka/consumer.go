package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokerURL, topic, groupID string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerURL},
		Topic:   topic,
		GroupID: groupID,
	})
	return &Consumer{reader: r}
}

func (c *Consumer) ReadMessage(ctx context.Context) (key, value []byte, err error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		log.Printf("failed to read message: %v", err)
		return nil, nil, err
	}

	return msg.Key, msg.Value, nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
