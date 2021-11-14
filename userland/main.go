package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/api/server"
	"github.com/werdna521/userland/db"
	"github.com/werdna521/userland/mailer"
	"github.com/werdna521/userland/producer"
)

func main() {
	// docker won't restart on failure if program is not running for at least 10
	// seconds
	fmt.Println("waiting")
	<-time.After(10 * time.Second)
	fmt.Println("waited")

	serverConfig := server.Config{
		Port: os.Getenv("API_PORT"),
	}
	postgresConfig := db.PostgresConfig{
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Addr:     os.Getenv("POSTGRES_ADDR"),
		Database: os.Getenv("POSTGRES_DB"),
	}
	redisConfig := db.RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}
	mailerConfig := mailer.Config{
		SenderName:  os.Getenv("SENDINBLUE_SENDER_NAME"),
		SenderEmail: os.Getenv("SENDINBLUE_SENDER_EMAIL"),
		APIKey:      os.Getenv("SENDINBLUE_API_KEY"),
	}
	producerConfig := producer.ProducerConfig{
		BootstrapServers: os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
	}

	log.Info().Msg("get connection to postgres")
	postgresConn, err := db.NewPosgresConn(postgresConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to postgres")
		panic("failed to connect to postgres")
	}

	log.Info().Msg("get connection to redis")
	redisConn, err := db.NewRedisConn(redisConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to redis")
		panic("failed to connect to redis")
	}

	log.Info().Msg("creating a kafka producer")
	p, err := producer.NewKafkaProducer(producerConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to create a kafka producer")
		panic("failed to create a kafka producer")
	}

	dataSource := &server.DataSource{
		Postgres: postgresConn,
		Redis:    redisConn,
	}

	mailer := mailer.NewBaseMailer(mailerConfig)

	log.Info().Msg("starting api server")
	server := server.NewServer(serverConfig, mailer, dataSource, p)
	server.Start()
}
