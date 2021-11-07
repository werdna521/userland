package db

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisConn(config RedisConfig) (*redis.Client, error) {
	log.Info().Msg("connecting to redis")
	rdb := redis.NewClient(getRedisOptions(config))

	log.Info().Msg("ping redis to check connection")
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		log.Error().Err(err).Stack().Msg("redis is not responding")
		return nil, err
	}

	return rdb, nil
}

func getRedisOptions(config RedisConfig) *redis.Options {
	return &redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	}
}
