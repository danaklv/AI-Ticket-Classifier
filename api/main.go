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
	db.AutoMigrate(&models.Ticket{}, &models.Outbox{})

	if err != nil {
		log.Fatal(err)
	}

	producer := kafka.NewProducer(cfg.KafkaBroker, cfg.KafkaTopic)
	defer producer.Close()

	a := app.NewApp(db, producer)
	a.Run()

}
