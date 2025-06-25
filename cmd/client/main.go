package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/qol"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"log"
)

func main() {

	fmt.Println("Starting Peril client...")
	conn := pubsub.GetConnection()
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Client welcome error: %v", err)
	}

	queueName := routing.PauseKey + "." + username
	chConn, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.QueueTypeTransient)
	if err != nil {
		log.Fatalf("Failed to declare and bind queue: %v", err)
	}
	defer chConn.Close()

	qol.WaitForSignalToKill()

	fmt.Println("Program aborted! Connection closed.")
}
