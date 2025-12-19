package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	kafka "github.com/segmentio/kafka-go"
)

var (
	url   = "localhost:19092"
	topic = "usernames"
)

type kafkaObj struct {
	URL   string
	Topic string
}

type Message struct {
	Msg       string
	Partition string
}

func main() {
	// KafkaQueue()
	// KafkaPubSub()
}

func KafkaQueue() {
	var wg sync.WaitGroup
	consumerGrp := "grp1"
	c := "c"

	k := kafkaObj{
		URL:   url,
		Topic: topic,
	}
	if err := k.write(context.Background(),
		[]Message{
			{
				Msg:       "message-1",
				Partition: "p1",
			},
			{
				Msg:       "message-2",
				Partition: "p2",
			},
			{
				Msg:       "message-3",
				Partition: "p3",
			},
		}); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go k.consume(consumerGrp, fmt.Sprintf("%v-%v", c, i), &wg)
	}
	wg.Wait()

}

func KafkaPubSub() {
	var wg sync.WaitGroup
	consumerGrp := "grp"
	c := "c"

	k := kafkaObj{
		URL:   url,
		Topic: topic,
	}
	if err := k.write(context.Background(),
		[]Message{
			{
				Msg:       "message-1",
				Partition: "p1",
			},
			{
				Msg:       "message-2",
				Partition: "p1",
			},
			{
				Msg:       "message-3",
				Partition: "p1",
			},
		}); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go k.consume(fmt.Sprintf("%v-%v", consumerGrp, i), fmt.Sprintf("%v-%v", c, i), &wg)
	}
	wg.Wait()

}

func (k *kafkaObj) write(ctx context.Context, msgs []Message) error {
	w := kafka.NewWriter(
		kafka.WriterConfig{

			Brokers:  []string{url},
			Topic:    k.Topic,
			Balancer: &kafka.Hash{},
		},
	)

	for _, m := range msgs {
		fmt.Printf("writing %v to partition %v\n", m.Msg, m.Partition)
		err := w.WriteMessages(ctx, kafka.Message{
			Value: []byte(m.Msg),
			Key:   []byte(m.Partition),
		})
		if err != nil {
			return err
		}
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	return nil
}

func (k *kafkaObj) consume(consumerGrp, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{k.URL},
		GroupID:  consumerGrp,
		Topic:    k.Topic,
		MaxBytes: 10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
			break
		}
		fmt.Printf("message at consumer/topic/partition/offset %s/%v/%v/%v: %s = %s\n", name, m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
