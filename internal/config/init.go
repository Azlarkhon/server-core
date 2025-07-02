package config

import (
	"os"
	"strings"
)

var (
	KafkaBrokers    []string
	TopicsToSend    []string
	MatchResults    string
	UsedItems string
)

func init() {

	// Kafka brokers
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	KafkaBrokers = strings.Split(kafkaBrokers, ",")

	// Topics to send
	MatchResults = os.Getenv("MATCH_RESULTS")
	UsedItems = os.Getenv("USED_ITEMS")

	TopicsToSend = []string{MatchResults, UsedItems}
}
