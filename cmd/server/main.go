package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	serverUrl := "amqp://guest:guest@localhost:5672/"
	serverConn, err := amqp.Dial(serverUrl)
	if err != nil {
		log.Fatalf("Error connecting to server: %s\n", err.Error())
	}
	defer serverConn.Close()
	fmt.Println("Connected to Peril server")

	serverChan, err := serverConn.Channel()
	if err != nil {
		log.Fatalf("Error creating channel over server connection: %s\n", err.Error())
	}

	perilPlayingState := routing.PlayingState{IsPaused: true}
	err = pubsub.PublishJSON(serverChan, routing.ExchangePerilDirect, routing.PauseKey, perilPlayingState)
	if err != nil {
		log.Fatalf("Error pubishing pause state: %s\n", err.Error())
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	fmt.Println("Shutting down...")
	err = serverConn.Close()
	if err != nil {
		log.Fatalf("Could not properly shut down server: %s\n", err.Error())
	}
	fmt.Println("Server has successfully shut down.")
}
