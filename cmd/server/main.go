package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/qol"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	fmt.Println("Starting Peril server...")
	connLink := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connLink)
	if err != nil {
		log.Fatal("Error connecting to RabbitMQ!")
	}

	fmt.Println("Connection was successful!")

	chConn, err := conn.Channel()
	if err != nil {
		log.Fatal("Can't open new channel")
	}

	data := routing.PlayingState{
		IsPaused: true,
	}
	pubsub.PublishJSON(chConn, routing.ExchangePerilDirect, routing.PauseKey, data)
	qol.WaitForSignalToKill()

	fmt.Println("Program aborted! Connection closed.")

	defer conn.Close()
	defer chConn.Close()
}
