package main

import (
	"consumer/logger"
	"consumer/rabbitmq/consumer"
	"context"
	"os"
	"os/signal"
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
	go func() {
		if err := consumer.DirectExchangeNormalQueue(context, log); err != nil {
			log.Fatal(err.Error())
		}
	}()

	go func() {
		if err := consumer.TopicExchangeQuorumQueue(context, log); err != nil {
			log.Fatal(err.Error())
		}
	}()

	go func() {
		if err := consumer.FanoutExchangeNormalQueue(context, log); err != nil {
			log.Fatal(err.Error())
		}
	}()

	go func() {
		if err := consumer.HeaderExchangeNormalQueue(context, log); err != nil {
			log.Fatal(err.Error())
		}
	}()
	<-sigChan
	cancel()
}
