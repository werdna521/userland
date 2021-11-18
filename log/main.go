package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/log/consumer"
	"github.com/werdna521/userland/log/core"
	"github.com/werdna521/userland/log/db"
)

func main() {
	// docker won't restart on failure if program is not running for at least 10
	// seconds
	fmt.Println("waiting")
	<-time.After(10 * time.Second)
	fmt.Println("waited")

	consumerConfig := consumer.ConsumerConfig{
		BootstrapServers: os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
		GroupID:          os.Getenv("KAFKA_GROUP_ID"),
		AutoOffsetReset:  os.Getenv("KAFKA_AUTO_OFFSET_RESET"),
	}
	postgresConfig := db.PostgresConfig{
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Addr:     os.Getenv("POSTGRES_ADDR"),
		Database: os.Getenv("POSTGRES_DB"),
	}

	c, err := consumer.NewKafkaConsumer(consumerConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to create kafka consumer")
		panic("kafka consumer failed")
	}

	postgresConn, err := db.NewPosgresConn(postgresConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to create postgres connection")
		panic("postgres connection failed")
	}

	d := core.Datasource{
		Postgres: postgresConn,
	}

	logMicro := core.NewBaseLogMicro(c, d)
	logMicro.InitLogService()
}
