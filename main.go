package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/api/server"
	"github.com/werdna521/userland/repository/postgres"
)

func main() {
	serverConfig := server.Config{
		Port: os.Getenv("API_PORT"),
	}

	postgresConfig := postgres.Config{
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	}

	log.Info().Msg("get a new postgres connection")
	postgresConn, err := postgres.NewPosgresConn(postgresConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to create a postgres connection")
		return
	}

	dataSource := &server.DataSource{
		Postgres: postgresConn,
	}

	log.Info().Msg("starting api server")
	server := server.NewServer(serverConfig, dataSource)
	server.Start()
}
