package pubsub

import "context"

// PubSub is pubsub interface.
//
// See usage example in example folder.
type PubSub interface {
	// Publish message to specific topic/channel.
	// Data will be encoded first before publishing.
	Publish(ctx context.Context, topic string, message []byte) error
	// Subscribe to specific topic/channel.
	Subscribe(ctx context.Context, topic string, handlerFunc HandlerFunc) error
	// Close pubsub client connection.
	Close() error
}

// HandlerFunc is pubsub subscriber handler function.
type HandlerFunc func(ctx context.Context, message []byte)
