package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
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

	gameState := gamelogic.NewGameState(username)

	pauseQueueName := routing.PauseKey + "." + username
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		pauseQueueName,
		routing.PauseKey,
		pubsub.QueueTypeTransient,
		handlerPause(gameState))
	if err != nil {
		log.Fatalf("Failed to subscribe to queue: %v", err)
	}

	armyMovesQueueName := "army_moves" + "." + username
	armyMovesRoutingKey := "army_moves.*"
	err = pubsub.SubscribeJSON(
		conn,
		string(routing.ExchangePerilTopic),
		armyMovesQueueName,
		armyMovesRoutingKey,
		pubsub.QueueTypeTransient,
		handlerMove(gameState))
	if err != nil {
		log.Fatalf("Failed to subscribe to '%v' queue: %v", armyMovesQueueName, err)
	}

	warQueueName := "war"
	warRoutingKey := "war.*"
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		warQueueName,
		warRoutingKey,
		pubsub.QueueTypeDurable,
		handleWar(gameState))
	if err != nil {
		log.Fatalf("Failed to subscribe to '%v' queue: %v", warQueueName, err)
	}
	playGame(gameState, username)
}

func playGame(gameState *gamelogic.GameState, username string) {
	for isOver := false; !isOver; {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch strings.TrimSpace(strings.ToLower(input[0])) {
		// spawn <location> <rank>
		case "spawn":
			err := gameState.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
			}
		// move <destination> <unit-id>
		case "move":
			armyMovesChannel := pubsub.GetChannel()
			armyMove, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Println(err)
			}
			err = pubsub.PublishJSON(armyMovesChannel, routing.ExchangePerilTopic, "army_moves."+username, armyMove)
			if err != nil {
				fmt.Printf("Error during moving army occured: %v", err)
			}
			fmt.Println("Your move was published!")
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
