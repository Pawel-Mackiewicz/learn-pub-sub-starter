package pubsub

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	connLink = "amqp://guest:guest@localhost:5672/"
)

var connInstance *amqp.Connection

func GetConnection() *amqp.Connection {
	if connInstance != nil {
		return connInstance
	}
	var err error
	connInstance, err = amqp.Dial(connLink)
	if err != nil {
		log.Fatalf("error connecting to RabbitMQ: %v", err)
	}
	fmt.Println("Connection was successful!")
	return connInstance
}

func GetChannel() *amqp.Channel {
	chConn, err := GetConnection().Channel()
	if err != nil {
		log.Fatalf("Can't open new channel: %v", err)
	}
	fmt.Println("New channel openned")
	return chConn
}
