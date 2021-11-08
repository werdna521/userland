package core

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
)

const userLoginTopic = "user_login"

type LogMicro interface {
	InitLogService()
}

type BaseLogMicro struct {
	c *kafka.Consumer
}

func NewBaseLogMicro(c *kafka.Consumer) *BaseLogMicro {
	return &BaseLogMicro{
		c: c,
	}
}

func (lm *BaseLogMicro) InitLogService() {
	topics := []string{userLoginTopic}
	lm.c.SubscribeTopics(topics, nil)

	exitChan := make(chan bool)

	for {
		select {
		case <-exitChan:
			log.Info().Msg("stopping log micro")
			return
		default:
			msg, err := lm.c.ReadMessage(-1)
			if err != nil {
				log.Error().Err(err).Msgf("consumer error: %v", err)
				exitChan <- true
				continue
			}

			switch *msg.TopicPartition.Topic {
			case userLoginTopic:
				log.Info().Msgf("message received: %s", string(msg.Value))
			default:
				log.Error().Msgf("unknown topic: %s", *msg.TopicPartition.Topic)
			}
		}
	}
}
