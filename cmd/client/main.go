package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/qol"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {

	fmt.Println("Starting Peril client...")
	connLink := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connLink)
	if err != nil {
		log.Fatal("Error connecting to RabbitMQ!")
	}
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Errorf("error: %v", err)
	}
	queueName := routing.PauseKey + "." + username
	pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.QueueTypeTransient)

	qol.WaitForSignalToKill()

	fmt.Println("Program aborted! Connection closed.")
	defer conn.Close()
}
