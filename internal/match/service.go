package match

import (
	"lesta-battleship/server-core/internal/config"
	"lesta-battleship/server-core/internal/infra/kafka"
)

type MatchService interface {
	HandleMatchResult(result MatchResult) error
	HandleUsedItem(item Item) error
}

type MatchEventPublisher struct {
	kafkaProducer kafka.KafkaProducer
}

func NewMatchEventPublisher(producer kafka.KafkaProducer) *MatchEventPublisher {
	return &MatchEventPublisher{kafkaProducer: producer}
}

func (s *MatchEventPublisher) HandleMatchResult(result MatchResult) error {
	return s.kafkaProducer.Send(config.MatchResults, result)
}

func (s *MatchEventPublisher) HandleUsedItem(item Item) error {
	return s.kafkaProducer.Send(config.UsedItems, item)
}
