package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Conn string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USEr")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	cfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

		return &Config{
			Conn: cfg,
		}
}
