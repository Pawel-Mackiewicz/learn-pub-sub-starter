package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

const (
	connLink = "amqp://guest:guest@localhost:5672/"
)

func GetConnection() *amqp.Connection {
	conn, err := amqp.Dial(connLink)
	if err != nil {
		log.Fatalf("error connecting to RabbitMQ: %v", err)
	}
	return conn
}
