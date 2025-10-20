package main

import (
	"classifier/internal/app"
	"classifier/internal/config"
	"classifier/internal/database"
	"classifier/internal/models"
	"log"
)

func main() {

	cfg := config.Load()

	db, err := database.Connect(cfg.Conn)
	db.AutoMigrate(&models.Ticket{})

	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(db)
	a.Run()

}
