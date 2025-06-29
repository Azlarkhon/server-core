package match

import (
	"log"

	"lesta-battleship/server-core/internal/kafka"
)

const (
	MatchResults    = "match-results"
	InventoryEvents = "inventory-events"
)

type MatchService interface {
	HandleMatchResult(result MatchResult) error
	HandleInventoryEvent(event InventoryEvent) error
}

type Service struct {
	kafkaProducer kafka.KafkaProducer
}

func NewService(producer kafka.KafkaProducer) *Service {
	return &Service{kafkaProducer: producer}
}

func (s *Service) HandleMatchResult(result MatchResult) error {
	if err := s.kafkaProducer.Send(MatchResults, result); err != nil {
		return err
	}

	log.Printf("Processed match result for match ID %d", result.MatchID)
	return nil
}

func (s *Service) HandleInventoryEvent(event InventoryEvent) error {
	if err := s.kafkaProducer.Send(InventoryEvents, event); err != nil {
		return err
	}

	log.Printf("Processed inventory event for player %d", event.PlayerID)
	return nil
}
