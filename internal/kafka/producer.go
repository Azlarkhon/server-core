package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type KafkaProducer interface {
	Send(topic string, message any) error
	Close() error
	CreateTopic(topic string, numPartitions int32, replicationFactor int16) error
}

type Producer struct {
	producer sarama.AsyncProducer
	admin    sarama.ClusterAdmin
	brokers  []string
}

func NewProducer(brokers []string, topics []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		asyncProducer.Close()
		return nil, err
	}

	producer := &Producer{
		producer: asyncProducer,
		admin:    admin,
		brokers:  brokers,
	}

	for _, topic := range topics {
		if err := producer.CreateTopic(topic, 3, 1); err != nil {
			log.Printf("Create topic failed: %v", err)
		}
	}

	go producer.listenForSuccess()
	go producer.listenForErrors()

	return producer, nil
}

func (p *Producer) listenForSuccess() {
	for success := range p.producer.Successes() {
		log.Printf("Message successfully sent to topic %s (partition %d, offset %d)",
			success.Topic, success.Partition, success.Offset)
	}
}

func (p *Producer) listenForErrors() {
	for err := range p.producer.Errors() {
		log.Printf("Failed to send message: %v", err)
	}
}

func (p *Producer) Send(topic string, message any) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return err
	}

	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msgBytes),
	}

	return nil
}

func (p *Producer) Close() error {
	err := p.producer.Close()
	if err != nil {
		return err
	}
	return p.admin.Close()
}

func (p *Producer) CreateTopic(topic string, numPartitions int32, replicationFactor int16) error {
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	err := p.admin.CreateTopic(topic, topicDetail, false)
	if err != nil {
		if terr, ok := err.(*sarama.TopicError); ok && terr.Err == sarama.ErrTopicAlreadyExists {
			log.Printf("Topic %s already exists", topic)
			return nil
		}
		log.Printf("Failed to create topic %s: %v", topic, err)
		return err
	}

	log.Printf("Successfully created topic %s with %d partitions and replication factor %d",
		topic, numPartitions, replicationFactor)
	return nil
}
