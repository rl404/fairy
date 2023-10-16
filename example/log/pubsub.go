package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/pubsub"
	"github.com/rl404/fairy/pubsub/rabbitmq"
)

type sampleData struct {
	Field1 string
	Field2 int
}

// Example use of rabbitmq pubsub middleware with log.
func pubsubWithLog(l log.Logger) {
	var client pubsub.PubSub
	var err error

	topic := "topic"

	// Init pubsub.
	client, err = rabbitmq.New("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(err)
	}

	// Don't forget to close.
	defer client.Close()

	// Add logger middleware.
	client.Use(log.PubSubMiddlewareWithLog(l, log.PubSubMiddlewareConfig{
		Topic:   topic,
		Payload: true,
		Error:   true,
	}))

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

	// Subscribe to specific topic/channel.
	if err := client.Subscribe(context.Background(), topic, handler); err != nil {
		panic(err)
	}

	// Publish data to specific topic/channel.
	if err := client.Publish(context.Background(), topic, jsonData); err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
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

	if err := sampleErr(ctx); err != nil {
		// Let's also test the error stack trace feature.
		stack.Wrap(ctx, errors.New("sample error"), err)
	}
}
