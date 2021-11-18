package core

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/log/repository"
	"github.com/werdna521/userland/log/repository/postgres"
)

const userLoginTopic = "user_login"

type LogMicro interface {
	InitLogService()
}

type BaseLogMicro struct {
	c *kafka.Consumer
	Datasource
	repositories *repositories
}
type Datasource struct {
	Postgres *sql.DB
}

type repositories struct {
	lr postgres.LogRepository
}

type userLoginMessage struct {
	UserID   string `json:"user_id"`
	RemoteIP string `json:"remote_ip"`
}

func NewBaseLogMicro(c *kafka.Consumer, d Datasource) *BaseLogMicro {
	return &BaseLogMicro{
		c:          c,
		Datasource: d,
	}
}

func (lm *BaseLogMicro) InitLogService() {
	ctx := context.Background()
	lm.initRepositories(ctx)

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
				userLoginMessage := &userLoginMessage{}
				log.Info().Msg("parsing json message")
				err := json.Unmarshal(msg.Value, userLoginMessage)
				if err != nil {
					log.Error().Err(err).Msg("fail to parse json message")
					// fail silently for now
					continue
				}

				auditLog := &repository.AuditLog{
					UserID:    userLoginMessage.UserID,
					RemoteIP:  userLoginMessage.RemoteIP,
					AuditType: userLoginTopic,
				}
				log.Info().Msg("saving user login log")
				_, err = lm.repositories.lr.CreateAuditLog(ctx, auditLog)
				if err != nil {
					log.Error().Err(err).Msg("fail to save user login log")
					// fail silently for now
					continue
				}

			default:
				log.Error().Msgf("unknown topic: %s", *msg.TopicPartition.Topic)
			}
		}
	}
}

func (lm *BaseLogMicro) initRepositories(ctx context.Context) {
	lr := postgres.NewBaseLogRepository(lm.Datasource.Postgres)
	lr.PrepareStatements(ctx)

	lm.repositories = &repositories{
		lr: lr,
	}
}
