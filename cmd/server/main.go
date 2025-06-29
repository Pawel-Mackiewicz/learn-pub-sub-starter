package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	logsQueueName = routing.GameLogSlug
	logsKeyName   = "game_logs.*"
	logsQueueType = pubsub.QueueTypeDurable
)

func main() {
	fmt.Println("Starting Peril server...")
	conn := pubsub.GetConnection()
	//nolint
	defer conn.Close()

	fmt.Println("Connection was successful!")

	pauseCh := pubsub.GetChannel()
	//nolint
	defer pauseCh.Close()

	_, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, logsQueueName, logsKeyName, logsQueueType)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	pauseOrResumeGame(pauseCh)
}

func pauseOrResumeGame(chConn *amqp.Channel) {
	gamelogic.PrintServerHelp()

	for quit := false; !quit; {
		words := gamelogic.GetInput()
		var pause bool
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Sending 'pause' message")
			pause = true
		case "resume":
			fmt.Println("Sending 'resume' message")
			pause = false
		case "quit":
			fmt.Println("Bye!")
			quit = true
		case "help":
			gamelogic.PrintServerHelp()
		default:
			fmt.Println("I don't understand this command :(")
		}
		publishPauseOrResume(chConn, pause)
	}
}

func publishPauseOrResume(chConn *amqp.Channel, pause bool) {
	data := routing.PlayingState{
		IsPaused: pause,
	}
	err := pubsub.PublishJSON(chConn, routing.ExchangePerilDirect, routing.PauseKey, data)
	if err != nil {
		fmt.Printf("Error publishing JSON: %v", err)
	}
}
