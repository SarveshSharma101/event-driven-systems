package publisher

import (
	"context"
	"os"
	"publisher/logger"
	"publisher/rabbitmq"
	"publisher/utils"

	"github.com/rabbitmq/amqp091-go"
)

var message string = "This is a message for the consumer"

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
		q1 := rabbitmq.Queue{
			Name:       "normalQ1",
			RoutingKey: "q1",
		}
		q2 := rabbitmq.Queue{
			Name:       "normalQ2",
			RoutingKey: "q2",
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

		for _, q := range rbmq.Q {
			if err := ch.PublishWithContext(
				ctx,
				rbmq.Exchange.Name,
				q.RoutingKey,
				false,
				false,
				amqp091.Publishing{
					ContentType: "text/plain",
					Body:        []byte(message),
				},
			); err != nil {

				log.Error("error publishing message to queue %v", q.Name)
				return err
			}
		}
	}
	return nil
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

		for _, q := range rbmq.Q {
			if err := ch.PublishWithContext(
				ctx,
				rbmq.Exchange.Name,
				q.Topic,
				false,
				false,
				amqp091.Publishing{
					ContentType: "text/plain",
					Body:        []byte(message),
				},
			); err != nil {

				log.Error("error publishing message to queue %v", q.Name)
				return err
			}
		}
	}
	return nil
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

		if err := ch.PublishWithContext(
			ctx,
			rbmq.Exchange.Name,
			rk,
			false,
			false,
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
			},
		); err != nil {

			log.Error("error publishing message to rk %v", rk)
			return err
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

		if err := ch.PublishWithContext(
			ctx,
			rbmq.Exchange.Name,
			"",
			false,
			false,
			amqp091.Publishing{
				ContentType: "text/plain",
				Body:        []byte(message),
				Headers: amqp091.Table{
					"a": "b",
					"c": "d",
				},
			},
		); err != nil {

			log.Error("error publishing message")
			return err
		}

	}
	return nil
}
