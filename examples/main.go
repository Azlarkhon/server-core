package main

import (
	"log"

	"lesta-battleship/server-core/internal/config"
	"lesta-battleship/server-core/internal/kafka"
	"lesta-battleship/server-core/internal/match"
	"lesta-battleship/server-core/internal/sample"
)

func main() {
	producer, err := kafka.NewProducer(config.KafkaBrokers, config.TopicsToSend)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	service := match.NewService(producer)
	handler := match.NewMatchHandler(service)

	// тест
	sample.RunSampleMatch(handler)
}
