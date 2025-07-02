package match

import (
	"lesta-battleship/server-core/internal/config"
	"lesta-battleship/server-core/internal/kafka"
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
	if err := s.kafkaProducer.Send(config.MatchResults, result); err != nil {
		return err
	}

	return nil
}

func (s *Service) HandleInventoryEvent(event InventoryEvent) error {
	if err := s.kafkaProducer.Send(config.InventoryEvents, event); err != nil {
		return err
	}

	return nil
}
