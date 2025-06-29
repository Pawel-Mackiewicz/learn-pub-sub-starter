package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

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
			makeWarCh := pubsub.GetChannel()
			rw := gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.Player,
			}
			warKey := routing.WarRecognitionsPrefix + "." + gs.GetUsername()
			err := pubsub.PublishJSON(makeWarCh, routing.ExchangePerilTopic, warKey, rw)
			if err != nil {
				fmt.Printf("Error during Publishing War!: %v", err)
				fmt.Println("Re-doing!")
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}
		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}

func handleWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		var logMsg string
		warOutcome, winner, loser := gs.HandleWar(rw)
		switch warOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeOpponentWon:
			fallthrough
		case gamelogic.WarOutcomeYouWon:
			logMsg = fmt.Sprintf("%v won a war against %v", winner, loser)
			fallthrough
		case gamelogic.WarOutcomeDraw:
			if logMsg == "" {
				logMsg = fmt.Sprintf("A war between %v and %v resulted in a draw", winner, loser)
			}
			key := routing.GameLogSlug + "." + rw.Attacker.Username
			err := pubsub.PublishGob(routing.ExchangePerilTopic, key, logMsg)
			if err != nil {
				fmt.Printf("error during publishing logs: %v", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			fmt.Println("Unexpected War Outcome")
			return pubsub.NackDiscard
		}
	}
}
