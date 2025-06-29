package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/rabbitmq/amqp091-go"
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

	pauseChan, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		pauseQueueName,
		routing.PauseKey,
		pubsub.QueueTypeTransient)
	if err != nil {
		log.Fatalf("Failed to declare and bind queue: %v", err)
	}
	//nolint
	defer pauseChan.Close()

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
	armyMovesChannel, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		armyMovesQueueName,
		armyMovesRoutingKey,
		pubsub.QueueTypeTransient)
	if err != nil {
		log.Fatalf("Failed to declare and bind '%v' queue: %v", armyMovesQueueName, err)
	}
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
	playGame(gameState, username, armyMovesChannel)
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		moveOutcome := gs.HandleMove(move)
		switch moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			return pubsub.Ack
		}
		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}

func playGame(gameState *gamelogic.GameState, username string, armyMovesChannel *amqp091.Channel) {
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
