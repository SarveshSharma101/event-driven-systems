package main

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

var (
	queueName   = "quorum-q"
	exhangeName = "quorum-exchange"
	routingKey  = "rk"
)

func main() {
	url := "amqp://default_user_vFYatzYwNIxc2qtbjgW:jTbfAh4Diemew_JjmpjLJnzufmd21vYV@localhost:5672/"
	conn, err := amqp091.Dial(url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"quorum-exchange",
		string("direct"),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	_, err = ch.QueueDeclare(
		"quorum-q",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-queue-type": amqp091.QueueTypeQuorum,
		},
	)
	if err != nil {
		panic(err)
	}

	if err := ch.QueueBind(
		"quorum-q",
		"rk",
		"quorum-exchange",
		false,
		nil,
	); err != nil {
		panic(err)
	}

	if err := ch.PublishWithContext(
		context.Background(),
		"quorum-exchange",
		"rk",
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte("message"),
		},
	); err != nil {

		panic(err)
	}

	msgs, err := ch.Consume(
		"quorum-q",
		"consumer",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		fmt.Printf("message from queue: %v", string(msg.Body))
		msg.Ack(false)
	}

}
