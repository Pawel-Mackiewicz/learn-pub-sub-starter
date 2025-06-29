package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackDiscard
	NackRequeue
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	deliveries, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	//handle deliveries
	go func() {
		for msg := range deliveries {
			var body T
			err := json.Unmarshal(msg.Body, &body)
			if err != nil {
				fmt.Printf("Error during JSON decoding: %v", err)
			}
			switch handler(body) {
			case Ack:
				err = msg.Ack(false)
				if err != nil {
					fmt.Println("Error during message acknowledgement")
				}
				fmt.Println("Ack")
			case NackDiscard:
				err = msg.Nack(false, false)
				if err != nil {
					fmt.Println("Error during message acknowledgement")
				}
				fmt.Println("NackDiscard")
			case NackRequeue:
				err = msg.Nack(false, true)
				if err != nil {
					fmt.Println("Error during message acknowledgement")
				}
				fmt.Println("NackRequeue")
			}
		}
	}()

	return nil
}
