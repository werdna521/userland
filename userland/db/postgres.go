package db

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
)

type PostgresConfig struct {
	Username string
	Password string
	Addr     string
	Database string
}

func NewPosgresConn(config PostgresConfig) (*sql.DB, error) {
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

	log.Info().Msg("ping postgres to check connection")
	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("postgres is not responding")
		return nil, err
	}

	return db, nil
}

func getConnURL(config PostgresConfig) string {
	// here, we're using the container name postgres to link to the postgres URL
	// (ordinarily should be hostname:port)
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		config.Username,
		config.Password,
		config.Addr,
		config.Database,
	)
}
