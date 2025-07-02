package match

import (
	"lesta-battleship/server-core/internal/config"
	"lesta-battleship/server-core/internal/kafka"
)

type MatchService interface {
	HandleMatchResult(result MatchResult) error
	HandleUsedItem(event Item) error
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

func (s *Service) HandleUsedItem(item Item) error {
	if err := s.kafkaProducer.Send(config.UsedItems, item); err != nil {
		return err
	}

	return nil
}
