package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	amqpConnectionString := "amqp://guest:guest@localhost:5672/"
	mqConn, err := amqp.Dial(amqpConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer mqConn.Close()

	fmt.Println("Connection is successful")
	fmt.Println("Starting Peril server...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan

	fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
}
