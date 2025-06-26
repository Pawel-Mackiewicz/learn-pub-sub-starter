package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/qol"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func main() {

	fmt.Println("Starting Peril client...")
	conn := pubsub.GetConnection()
	//nolint
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

	gameState := gamelogic.NewGameState(username)
	playGame(gameState)

	//nolint
	defer chConn.Close()

	qol.WaitForSignalToKill()

	fmt.Println("Program aborted! Connection closed.")
}

func playGame(gameState *gamelogic.GameState) {
	for isOver := false; !isOver; {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch strings.TrimSpace(strings.ToLower(input[0])) {
		//spawn <location> <rank>
		case "spawn":
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
			}
		//move <destination> <unit-id>
		case "move":
			_, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			isOver = true
		default:
			fmt.Println("I can't understand Your command! Try using `help` :)")
		}
	}
}
