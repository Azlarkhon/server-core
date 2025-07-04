package main

import (
	"log"

	"lesta-battleship/server-core/internal/config"
	"lesta-battleship/server-core/internal/infra/kafka"
	"lesta-battleship/server-core/internal/match"
	"lesta-battleship/server-core/internal/sample"
)

func main() {
	producer, err := kafka.NewProducer(config.KafkaBrokers, config.TopicsToSend)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	eventPublisher := match.NewMatchEventPublisher(producer)
	handler := match.NewMatchHandler(eventPublisher)

	// тест
	sample.RunSampleMatch(handler)
}
