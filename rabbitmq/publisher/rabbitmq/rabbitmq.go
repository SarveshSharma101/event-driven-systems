package rabbitmq

import (
	"publisher/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	Name       string
	Args       amqp.Table
	RoutingKey string
	Topic      string
}

type Exchange struct {
	Name string
	Type utils.ExchangeType
}

type Rabbitmq struct {
	Exchange   Exchange
	Q          []Queue
	Connection *amqp.Connection
	Url        string
}

type Message struct {
	Body string
}

func (r *Rabbitmq) Init() error {
	conn, err := amqp.Dial(r.Url)
	if err != nil {
		return err
	}
	r.Connection = conn
	return nil
}

func (r Rabbitmq) Bind(ch *amqp.Channel) error {
	for _, q := range r.Q {
		if err := ch.QueueBind(
			q.Name,
			q.RoutingKey,
			r.Exchange.Name,
			false,
			q.Args,
		); err != nil {
			return err
		}
	}
	return nil
}

func (e Exchange) CreateExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		e.Name,
		string(e.Type),
		true,
		false,
		false,
		false,
		nil,
	)
}

func (q Queue) CreateQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		q.Name,
		true,
		false,
		false,
		false,
		q.Args,
	)
}
