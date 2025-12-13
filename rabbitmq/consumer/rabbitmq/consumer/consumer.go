package consumer

import (
	"consumer/logger"
	"consumer/rabbitmq"
	"consumer/utils"
	"context"
	"fmt"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func DirectExchangeNormalQueue(ctx context.Context, log *logger.Logger) error {

	select {
	case <-ctx.Done():
		log.Info("context cancelled")
		return nil
	default:
		log.Info("Starting flow for direct exchange and normal queue")
		url := os.Getenv(utils.URL)
		directEx := rabbitmq.Exchange{
			Name: "DirectExchange-normalQueue",
			Type: utils.DIRECT,
		}

		deadEx := rabbitmq.Exchange{
			Name: "DeadExchange",
			Type: utils.DIRECT,
		}

		dlq := rabbitmq.Queue{
			Name:       "dead-q",
			RoutingKey: "dq",
		}

		q1 := rabbitmq.Queue{
			Name:       "normalQ1",
			RoutingKey: "q1",
		}
		q2 := rabbitmq.Queue{
			Name:       "normalQ2",
			RoutingKey: "q2",
			Args: amqp091.Table{
				"x-dead-letter-exchange":    deadEx.Name,
				"x-dead-letter-routing-key": dlq.RoutingKey,
			},
		}

		rbmq := rabbitmq.Rabbitmq{
			Exchange: directEx,
			Q:        []rabbitmq.Queue{q1, q2},
			Url:      url,
		}

		if err := rbmq.Init(); err != nil {
			log.Error("error initializing the rabbitmq connection on url %v", url)
			return err
		}
		defer rbmq.Connection.Close()
		ch, err := rbmq.Connection.Channel()
		if err != nil {
			log.Error("error initializing the channel")
			return err
		}

		defer ch.Close()

		if err := directEx.CreateExchange(ch); err != nil {
			log.Error("error creating exchange")
			return err
		}

		if err := deadEx.CreateExchange(ch); err != nil {
			log.Error("error creating exchange")
			return err
		}

		for _, q := range rbmq.Q {
			if _, err := q.CreateQueue(ch); err != nil {
				log.Error("error creating queue")
				return err
			}
		}

		if _, err := dlq.CreateQueue(ch); err != nil {
			log.Error("error creating exchange")
			return err
		}

		if err := rbmq.Bind(ch); err != nil {
			log.Error("error initializing the rabbitmq connection on url %v", url)
			return err
		}

		msgs1, err := ch.Consume(
			q1.Name,
			"q1-consumer",
			false, // manual ACK
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		msgs2, err := ch.Consume(
			q2.Name,
			"q2-consumer",
			false, // manual ACK
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		for {
			select {
			case m := <-msgs1:
				fmt.Printf("message from q1: %v", m)
				m.Ack(false)

			case m := <-msgs2:
				fmt.Printf("message from q1: %v", m)
				m.Nack(false, false)
			}

		}
	}
}

func TopicExchangeQuorumQueue(ctx context.Context, log *logger.Logger) error {

	select {
	case <-ctx.Done():
		log.Info("context cancelled")
		return nil
	default:
		log.Info("Starting flow for topic exchange and normal queue")
		url := os.Getenv(utils.URL)
		directEx := rabbitmq.Exchange{
			Name: "TopicExchange-normalQueue",
			Type: utils.TOPIC,
		}
		q1 := rabbitmq.Queue{
			Name:       "orderQ",
			RoutingKey: "*.Order.*",
			Topic:      "a.Order.b",
			Args: amqp091.Table{
				"x-queue-type": amqp091.QueueTypeQuorum,
			},
		}
		q2 := rabbitmq.Queue{
			Name:       "payQ",
			RoutingKey: "pay.#",
			Topic:      "pay.from.atm",
		}

		rbmq := rabbitmq.Rabbitmq{
			Exchange: directEx,
			Q:        []rabbitmq.Queue{q1, q2},
			Url:      url,
		}

		if err := rbmq.Init(); err != nil {
			log.Error("error initializing the rabbitmq connection on url %v", url)
			return err
		}
		defer rbmq.Connection.Close()
		ch, err := rbmq.Connection.Channel()
		if err != nil {
			log.Error("error initializing the channel")
			return err
		}

		defer ch.Close()

		if err := directEx.CreateExchange(ch); err != nil {
			log.Error("error creating exchange")
			return err
		}

		for _, q := range rbmq.Q {
			if _, err := q.CreateQueue(ch); err != nil {
				log.Error("error creating queue")
				return err
			}
		}

		if err := rbmq.Bind(ch); err != nil {
			log.Error("error initializing the rabbitmq connection on url %v", url)
			return err
		}

		msgs1, err := ch.Consume(
			q1.Name,
			"q1-consumer",
			false, // manual ACK
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		msgs2, err := ch.Consume(
			q2.Name,
			"q2-consumer",
			false, // manual ACK
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		for {
			select {
			case m := <-msgs1:
				fmt.Printf("message from q1: %v", m)
				m.Ack(false)

			case m := <-msgs2:
				fmt.Printf("message from q1: %v", m)
				m.Nack(false, false)
			}

		}
	}
}

func FanoutExchangeNormalQueue(ctx context.Context, log *logger.Logger) error {

	select {
	case <-ctx.Done():
		log.Info("context cancelled")
		return nil
	default:
		log.Info("Starting flow for direct exchange and normal queue")
		url := os.Getenv(utils.URL)
		ex := rabbitmq.Exchange{
			Name: "FanoutExchange-normalQueue",
			Type: utils.FANOUT,
		}
		rk := "fanout"
		q1 := rabbitmq.Queue{
			Name:       "Q1",
			RoutingKey: rk,
		}
		q2 := rabbitmq.Queue{
			Name:       "Q2",
			RoutingKey: rk,
		}

		rbmq := rabbitmq.Rabbitmq{
			Exchange: ex,
			Q:        []rabbitmq.Queue{q1, q2},
			Url:      url,
		}

		if err := rbmq.Init(); err != nil {
			log.Error("error initializing the rabbitmq connection on url %v", url)
			return err
		}
		defer rbmq.Connection.Close()
		ch, err := rbmq.Connection.Channel()
		if err != nil {
			log.Error("error initializing the channel")
			return err
		}

		defer ch.Close()

		if err := ex.CreateExchange(ch); err != nil {
			log.Error("error creating exchange")
			return err
		}

		for _, q := range rbmq.Q {
			if _, err := q.CreateQueue(ch); err != nil {
				log.Error("error creating queue")
				return err
			}
		}

		if err := rbmq.Bind(ch); err != nil {
			log.Error("error initializing the rabbitmq connection on url %v", url)
			return err
		}

		msgs1, err := ch.Consume(
			q1.Name,
			"q1-consumer",
			false, // manual ACK
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		msgs2, err := ch.Consume(
			q2.Name,
			"q2-consumer",
			false, // manual ACK
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		for {
			select {
			case m := <-msgs1:
				fmt.Printf("message from q1: %v", m)
				m.Ack(false)

			case m := <-msgs2:
				fmt.Printf("message from q1: %v", m)
				m.Nack(false, false)
			}

		}

	}
	return nil
}

func HeaderExchangeNormalQueue(ctx context.Context, log *logger.Logger) error {

	select {
	case <-ctx.Done():
		log.Info("context cancelled")
		return nil
	default:
		log.Info("Starting flow for direct exchange and normal queue")
		url := os.Getenv(utils.URL)
		ex := rabbitmq.Exchange{
			Name: "HeaderExchange-normalQueue",
			Type: utils.HEADERS,
		}
		q1 := rabbitmq.Queue{
			Name: "HQ1",
			Args: amqp091.Table{
				"x-match": "all",
				"a":       "b",
				"c":       "d",
			},
		}
		q2 := rabbitmq.Queue{
			Name: "HQ2",
			Args: amqp091.Table{
				"x-match": "any",
				"a":       "b",
			},
		}

		rbmq := rabbitmq.Rabbitmq{
			Exchange: ex,
			Q:        []rabbitmq.Queue{q1, q2},
			Url:      url,
		}

		if err := rbmq.Init(); err != nil {
			log.Error("error initializing the rabbitmq connection on url %v", url)
			return err
		}
		defer rbmq.Connection.Close()
		ch, err := rbmq.Connection.Channel()
		if err != nil {
			log.Error("error initializing the channel")
			return err
		}

		defer ch.Close()

		if err := ex.CreateExchange(ch); err != nil {
			log.Error("error creating exchange")
			return err
		}

		for _, q := range rbmq.Q {
			if _, err := q.CreateQueue(ch); err != nil {
				log.Error("error creating queue")
				return err
			}
		}

		if err := rbmq.Bind(ch); err != nil {
			log.Error("error initializing the rabbitmq connection on url %v", url)
			return err
		}

		msgs1, err := ch.Consume(
			q1.Name,
			"q1-consumer",
			false, // manual ACK
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		msgs2, err := ch.Consume(
			q2.Name,
			"q2-consumer",
			false, // manual ACK
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		for {
			select {
			case m := <-msgs1:
				fmt.Printf("message from q1: %v", m)
				m.Ack(false)

			case m := <-msgs2:
				fmt.Printf("message from q1: %v", m)
				m.Nack(false, false)
			}

		}

	}
	return nil
}
