package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerConfig struct {
	BootstrapServers string
	GroupID          string
	AutoOffsetReset  string
}

func NewKafkaConsumer(config ConsumerConfig) (*kafka.Consumer, error) {
	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"group.id":          config.GroupID,
		"auto.offset.reset": config.AutoOffsetReset,
	})
}
