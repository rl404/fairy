package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type sampleData struct {
	Field1 string
	Field2 int
}

func main() {
	// Pubsub type.
	t := RabbitMQ

	// Init client.
	client, err := New(t, "amqp://guest:guest@localhost:5672", "")
	if err != nil {
		panic(err)
	}

	// Don't forget to close.
	defer client.Close()

	// Prepare goroutine channel that will stop when ctrl+c.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Subscribe to specific topic/channel.
	if err := client.Subscribe(context.Background(), "topic", handler); err != nil {
		panic(err)
	}

	<-sigChan
}

// handler to handle the incoming message.
func handler(ctx context.Context, data []byte) {
	// Convert to the original struct when you publish it.
	var sampleData sampleData
	if err := json.Unmarshal(data, &sampleData); err != nil {
		panic(err)
	}

	// Process the message.
	fmt.Println(sampleData.Field1, sampleData.Field2)
}
