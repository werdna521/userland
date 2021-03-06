package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/api/server"
	"github.com/werdna521/userland/db"
	"github.com/werdna521/userland/mailer"
)

func main() {
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

	log.Info().Msg("get connection to postgres")
	postgresConn, err := db.NewPosgresConn(postgresConfig)
	if err != nil {
		log.Error().Err(err).Stack().Msg("failed to connect to postgres")
		return
	}

	log.Info().Msg("get connection to redis")
	redisConn, err := db.NewRedisConn(redisConfig)
	if err != nil {
		log.Error().Err(err).Stack().Msg("failed to connect to redis")
		return
	}

	dataSource := &server.DataSource{
		Postgres: postgresConn,
		Redis:    redisConn,
	}

	mailer := mailer.NewBaseMailer(mailerConfig)

	log.Info().Msg("starting api server")
	server := server.NewServer(serverConfig, mailer, dataSource)
	server.Start()
}
