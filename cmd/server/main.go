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
	fmt.Println("Starting Peril server...")
	conn := pubsub.GetConnection()
	defer conn.Close()

	fmt.Println("Connection was successful!")

	chConn, err := conn.Channel()
	if err != nil {
		log.Fatalf("Can't open new channel: %v", err)
	}
	defer chConn.Close()

	gamelogic.PrintServerHelp()
	for {
		words := gamelogic.GetInput()
		if words == nil || len(words) == 0 {
			continue
		}
		if words[0] == "pause" {
			fmt.Println("Sending 'pause' message")
			pauseResumeGame(chConn, "pause")
		} else if words[0] == "resume" {
			fmt.Println("Sending 'resume' message")
			pauseResumeGame(chConn, "resume")
		} else if words[0] == "quit" {
			fmt.Println("Bye!")
			break
		} else {
			fmt.Println("I don't understand the command :(")
		}
	}
	qol.WaitForSignalToKill()
	fmt.Println("Program aborted! Connection closed.")
}

func pauseResumeGame(chConn *amqp.Channel, state string) {
	var pause bool
	if state == "pause" {
		pause = true
	} else if state == "resume" {
		pause = false
	} else {
		log.Fatal("You called 'pauseResumeGame()' wrong!")
	}
	data := routing.PlayingState{
		IsPaused: pause,
	}
	err := pubsub.PublishJSON(chConn, routing.ExchangePerilDirect, routing.PauseKey, data)
	if err != nil {
		fmt.Printf("Error publishing JSON: %v", err)
	}

}
