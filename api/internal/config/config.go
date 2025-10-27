package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Conn string
	KafkaBroker string
	KafkaProducerTopic string
	KafkaConsumerTopic string
	
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	kafkaBroker := os.Getenv("KAFKA_BROKER")
	kafkaConsumerTopic := os.Getenv("KAFKA_CONSUMER_TOPIC")
	kafkaProducerTopic := os.Getenv("KAFKA_PRODUCER_TOPIC")

	cfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return &Config{
		Conn: cfg,
		KafkaBroker: kafkaBroker, 
		KafkaProducerTopic: kafkaProducerTopic,
		KafkaConsumerTopic: kafkaConsumerTopic,
	}
}
