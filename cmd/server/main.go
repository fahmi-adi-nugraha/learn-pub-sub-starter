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
	fmt.Println("Starting Peril server...")

	serverUrl := "amqp://guest:guest@localhost:5672/"
	serverConn, err := amqp.Dial(serverUrl)
	if err != nil {
		log.Fatalf("Error connecting to server: %s\n", err)
	}
	defer serverConn.Close()
	fmt.Println("Connected to Peril server")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	fmt.Println("Shutting down...")
	err = serverConn.Close()
	if err != nil {
		log.Fatalf("Could not properly shut down server: %s\n", err)
	}
	fmt.Println("Server has successfully shut down.")
}
