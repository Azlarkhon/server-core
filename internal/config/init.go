package config

import (
	"os"
	"strings"
)

var (
	KafkaBrokers    []string
	TopicsToSend    []string
	MatchResults    string
	InventoryEvents string
)

func init() {

	// Kafka brokers
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	KafkaBrokers = strings.Split(kafkaBrokers, ",")

	// Topics to send
	MatchResults = os.Getenv("MATCH_RESULTS")
	InventoryEvents = os.Getenv("INVENTORY_EVENTS")

	TopicsToSend = []string{MatchResults, InventoryEvents}
}
