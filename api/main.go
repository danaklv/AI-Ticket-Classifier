package main

import (
	"classifier/internal/app"
	"classifier/internal/config"
	"classifier/internal/database"
	"classifier/internal/kafka"
	"classifier/internal/models"
	"log"
)

func main() {

	cfg := config.Load()

	db, err := database.Connect(cfg.Conn)

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Ticket{}, &models.Outbox{})

	producer := kafka.NewProducer(cfg.KafkaBroker, cfg.KafkaProducerTopic)

	defer producer.Close()

	consumer := kafka.NewConsumer(cfg.KafkaBroker, cfg.KafkaConsumerTopic, "go-result-consumer")

	defer consumer.Close()

	a := app.NewApp(db, producer, consumer)
	a.Run()

}
