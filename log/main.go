package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/werdna521/userland/log/consumer"
	"github.com/werdna521/userland/log/core"
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

	c, err := consumer.NewKafkaConsumer(consumerConfig)
	if err != nil {
		log.Error().Err(err).Msg("failed to create kafka consumer")
		panic("kafka consumer failed")
	}

	logMicro := core.NewBaseLogMicro(c)
	logMicro.InitLogService()
}
