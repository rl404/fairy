package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rl404/fairy"
)

type sampleData struct {
	Field1 string
	Field2 int
}

func main() {
	// Pubsub type.
	t := fairy.RabbitMQ

	// Init client.
	client, err := fairy.NewPubSub(t, "amqp://guest:guest@localhost:5672", "")
	if err != nil {
		panic(err)
	}

	// Don't forget to close.
	defer client.Close()

	// Subscribe to specific topic/channel.
	s, err := client.Subscribe("topic")
	if err != nil {
		panic(err)
	}

	// Need to convert to Channel interface first
	// so you can use function in the Channel interface.
	channel := s.(fairy.Channel)

	// Don't forget to close subscription.
	defer channel.Close()

	// Prepare a new or existing variable for
	// incoming message. Data type should be the
	// same as when publish the message.
	var newData sampleData

	// Read incomming message. Message will be decoded
	// to newData. Don't forget to use pointer.
	msgs, errChan := channel.Read(&newData)

	// Prepare goroutine channel that will stop when
	// ctrl+c.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		// Loop for waiting incoming message.
		for {
			select {
			// If message comes.
			case <-msgs:
				// Process the message.
				fmt.Println(newData.Field1, newData.Field2)

			// If error comes.
			case err = <-errChan:
				// Process the error.
				fmt.Println(err)
			}
		}
	}()

	<-sigChan
}
