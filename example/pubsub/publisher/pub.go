package main

import (
	"context"
	"encoding/json"
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

	// Sample data. Can be any type.
	data := sampleData{
		Field1: "a",
		Field2: 1,
	}

	// Convert to []byte.
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// Publish data to specific topic/channel.
	if err := client.Publish(context.Background(), "topic", jsonData); err != nil {
		panic(err)
	}
}
