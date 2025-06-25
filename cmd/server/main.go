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
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}
	defer conn.Close()

	fmt.Println("Connection was successful!")

	chConn, err := conn.Channel()
	if err != nil {
		log.Fatalf("Can't open new channel: %v", err)
	}
	defer chConn.Close()

	data := routing.PlayingState{
		IsPaused: true,
	}
	err = pubsub.PublishJSON(chConn, routing.ExchangePerilDirect, routing.PauseKey, data)
	if err != nil {
		fmt.Printf("Error publishing JSON: %v", err)
	}

	qol.WaitForSignalToKill()
	fmt.Println("Program aborted! Connection closed.")
}
