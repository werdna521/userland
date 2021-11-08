package producer

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

const (
	userLoginTopic = "user_login"
)

type LoginLog struct {
	UserID   string `json:"user_id"`
	RemoteIP string `json:"remote_ip"`
}

type LogProducer interface {
	ProduceLoginTopic(ll *LoginLog) error
}

type BaseLogProducer struct {
	kp *kafka.Producer
}

func NewBaseLogProducer(kp *kafka.Producer) *BaseLogProducer {
	return &BaseLogProducer{
		kp: kp,
	}
}

func (p *BaseLogProducer) ProduceLoginTopic(ll *LoginLog) error {
	log.Info().Msg("stringifying login log struct")
	value, err := json.Marshal(ll)
	if err != nil {
		log.Error().Err(err).Msg("failed to stringify login log struct")
		return err
	}

	errChan := make(chan error)

	go func() {
		defer close(errChan)
		for e := range p.kp.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					log.Error().Err(m.TopicPartition.Error).Msg("error producing to kafka")
					errChan <- m.TopicPartition.Error
					return
				}
				return
			}
		}
	}()

	topic := userLoginTopic
	msg := newKafkaMessage(&topic, value)
	p.kp.ProduceChannel() <- msg

	err = <-errChan
	if err != nil {
		log.Error().Err(err).Msg("error producing to kafka")
		return err
	}

	return nil
}
