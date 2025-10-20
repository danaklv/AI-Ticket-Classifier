package main

import (
	"classifier/internal/app"
	"classifier/internal/config"
	"classifier/internal/database"
	"log"
)

func main() {

	cfg := config.Load()

	db, err := database.Connect(cfg.Conn)

	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(db)
	a.Run()

}
