package producer

import "github.com/confluentinc/confluent-kafka-go/kafka"

type ProducerConfig struct {
	BootstrapServers string
}

func NewKafkaProducer(config ProducerConfig) (*kafka.Producer, error) {
	return kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
	})
}

func newKafkaMessage(topic *string, value []byte) *kafka.Message {
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}
}
