package main

import (
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

	// Sample data. Can be any type.
	data := sampleData{
		Field1: "a",
		Field2: 1,
	}

	// Publish data to specific topic/channel. Data will be encoded first.
	if err = client.Publish("topic", data); err != nil {
		panic(err)
	}
}
