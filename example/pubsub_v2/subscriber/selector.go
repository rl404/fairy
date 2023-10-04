// This file contains function to help you switch pubsub type easily.
//
// You can use this in your project if you want.
//
// You don't need to handle all pubsub types, just handle
// the ones you will use.
package main

import (
	"errors"

	pubsub "github.com/rl404/fairy/pubsub_v2"
	"github.com/rl404/fairy/pubsub_v2/google"
	"github.com/rl404/fairy/pubsub_v2/kafka"
	"github.com/rl404/fairy/pubsub_v2/nsq"
	"github.com/rl404/fairy/pubsub_v2/rabbitmq"
	"github.com/rl404/fairy/pubsub_v2/redis"
)

// PubsubType is type for pubsub.
type PubsubType int8

// Available types for pubsub.
const (
	Redis PubsubType = iota + 1
	RabbitMQ
	NSQ
	Google
	Kafka
)

// ErrInvalidPubsubType is error for invalid pubsub type.
var ErrInvalidPubsubType = errors.New("invalid pubsub type")

// New to create new pubsub client depends on the type.
func New(pubsubType PubsubType, address string, password string) (pubsub.PubSub, error) {
	switch pubsubType {
	case Redis:
		return redis.New(address, password)
	case RabbitMQ:
		return rabbitmq.New(address)
	case NSQ:
		return nsq.New(address)
	case Google:
		return google.New(address, password)
	case Kafka:
		return kafka.New(address)
	default:
		return nil, ErrInvalidPubsubType
	}
}
