package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Username string
	Password string
	Database string
}

func NewPosgresConn(config Config) (*sql.DB, error) {
	connURL := getConnURL(config)

	log.Info().Msg("parsing conn config from conn URL")
	connConfig, err := pgx.ParseConfig(connURL)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse the config")
		return nil, err
	}

	connString := stdlib.RegisterConnConfig(connConfig)

	log.Info().Msg("opening a connection to postgres")
	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Error().Err(err).Msg("failed to open postgres connection")
		return nil, err
	}

	log.Info().Msg("pinging postgres")
	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("failed to ping postgres")
		return nil, err
	}

	return db, nil
}

func getConnURL(config Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@postgres/%s",
		config.Username,
		config.Password,
		config.Database,
	)
}
