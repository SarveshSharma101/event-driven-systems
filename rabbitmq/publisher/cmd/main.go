package main

import (
	"context"
	"os"
	"os/signal"
	"publisher/logger"
	"publisher/rabbitmq/publisher"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	log := logger.Get()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file %v", err)
	}
	context, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGKILL)

	log.Info("Starting rabbitmq Publisher")
	if err := publisher.DirectExchangeNormalQueue(context, log); err != nil {
		log.Fatal(err.Error())
	}

	if err := publisher.TopicExchangeQuorumQueue(context, log); err != nil {
		log.Fatal(err.Error())
	}

	if err := publisher.FanoutExchangeNormalQueue(context, log); err != nil {
		log.Fatal(err.Error())
	}

	if err := publisher.HeaderExchangeNormalQueue(context, log); err != nil {
		log.Fatal(err.Error())
	}
	<-sigChan
	cancel()

}
