package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	k := fmt.Sprintf("%s.*", routing.GameLogSlug)
	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, k, pubsub.SimpleQueueDurable)
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	gamelogic.PrintServerHelp()

	for {
		cmds := gamelogic.GetInput()
		if len(cmds) == 0 {
			continue
		}

		cmd := cmds[0]
		switch cmd {
		case "pause":
			fmt.Printf("publishing %s message\n", cmd)
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v\n", err)
			}
			fmt.Printf("%s message published!\n", cmd)
		case "resume":
			fmt.Printf("publishing %s message\n", cmd)
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v\n", err)
			}
			fmt.Printf("%s message published!\n", cmd)
		case "quit":
			fmt.Printf("got %s word, exitting\n", cmd)
			return
		default:
			fmt.Printf("%s word is eligible\n", cmd)
		}
	}
}
