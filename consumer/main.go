package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

type MatchResult struct {
	WinnerID           int         `json:"winner_id"`
	LoserID            int         `json:"loser_id"`
	MatchID            int         `json:"match_id"`
	MatchDurationInSec int         `json:"match_duration_in_sec"`
	MatchDate          time.Time   `json:"match_date"`
	MatchType          string      `json:"match_type"`
	Experience         *Experience `json:"experience,omitempty"`
}

type Experience struct {
	WinnerGain int `json:"winner_gain,omitempty"`
	LoserGain  int `json:"loser_gain,omitempty"`
}

type InventoryEvent struct {
	PlayerID int `json:"player_id"`
	ItemID   int `json:"item_id"`
}

func main() {
	topics := []string{"match-results", "inventory-events"}
	brokers := []string{"kafka:9092"}

	consumer, err := ConnectConsumer(brokers)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	}()

	partitionConsumers := make([]sarama.PartitionConsumer, 0, len(topics))
	for _, topic := range topics {
		pc, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Fatalf("Failed to start consumer for topic %s: %v", topic, err)
		}
		defer pc.Close()
		partitionConsumers = append(partitionConsumers, pc)
	}

	log.Println("Consumer started successfully for topics:", topics)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	msgCount := 0
	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case <-sigchan:
				log.Println("Interrupt signal received, shutting down...")
				doneCh <- struct{}{}
				return

			default:
				for _, pc := range partitionConsumers {
					select {
					case err := <-pc.Errors():
						log.Printf("Error from consumer: %v", err)
					case msg := <-pc.Messages():
						msgCount++
						processMessage(msg)
					default:
						time.Sleep(100 * time.Millisecond)
					}
				}
			}
		}
	}()

	<-doneCh
	log.Printf("Processed %d messages before shutdown", msgCount)
}

func processMessage(msg *sarama.ConsumerMessage) {
	log.Printf("Received message from topic %s: %s", msg.Topic, string(msg.Value))

	switch msg.Topic {
	case "match-results":
		var result MatchResult
		if err := json.Unmarshal(msg.Value, &result); err != nil {
			log.Printf("Error parsing match result: %v", err)
			return
		}
		handleMatchResult(result)

	case "inventory-events":
		var event InventoryEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error parsing inventory event: %v", err)
			return
		}
		handleInventoryEvent(event)

	default:
		log.Printf("Unknown topic: %s", msg.Topic)
	}
}

func handleMatchResult(result MatchResult) {
	log.Printf("Processing match result - Winner: %d, Loser: %d", result.WinnerID, result.LoserID)
	log.Printf("Match duration: %d seconds, Date: %s", result.MatchDurationInSec, result.MatchDate)
	log.Printf("Experience - Winner gains: %d, Loser gains: %d", result.Experience.WinnerGain, result.Experience.LoserGain)
}

func handleInventoryEvent(event InventoryEvent) {
	log.Printf("Processing inventory event - player_id: %d, item_id: %d", event.PlayerID, event.ItemID)
}

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 10 * time.Second
	config.Net.WriteTimeout = 10 * time.Second
	config.Metadata.Retry.Max = 5
	config.Metadata.Retry.Backoff = 1 * time.Second

	var consumer sarama.Consumer
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		consumer, err = sarama.NewConsumer(brokers, config)
		if err == nil {
			return consumer, nil
		}
		log.Printf("Attempt %d/%d: Failed to connect to Kafka: %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to Kafka after %d attempts: %v", maxRetries, err)
}
