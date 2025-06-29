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
	conn, err := amqp.Dial(connLink)
	if err != nil {
		log.Fatalf("error connecting to RabbitMQ: %v", err)
	}
	fmt.Println("Connection was successful!")
	return conn
}

func GetChannel() *amqp.Channel {
	chConn, err := GetConnection().Channel()
	if err != nil {
		log.Fatalf("Can't open new channel: %v", err)
	}
	fmt.Println("Channel opened successfully!")
	return chConn
}
