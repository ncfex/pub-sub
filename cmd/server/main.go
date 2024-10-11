package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func gobDecoder[T any](data []byte) (T, error) {
	var target T
	decoder := gob.NewDecoder(bytes.NewBuffer(data))
	err := decoder.Decode(&target)
	return target, err
}

func jsonDecoder[T any](data []byte) (T, error) {
	var target T
	err := json.Unmarshal(data, &target)
	return target, err
}

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

	err = pubsub.SubscribeGOB(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
		handlerLog(),
		gobDecoder,
	)
	if err != nil {
		log.Fatalf("could not starting consuming logs: %v", err)
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
