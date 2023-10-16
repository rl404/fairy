// Package nop is no-operation pubsub.
package nop

import (
	"context"

	"github.com/rl404/fairy/pubsub"
)

// Client is no-operation pubsub client.
type Client struct{}

// New to create new no-operation pubsub client.
func New() *Client {
	return &Client{}
}

// Use to do nothing.
func (c *Client) Use(middlewares ...func(pubsub.HandlerFunc) pubsub.HandlerFunc) {}

// Publish to do nothing.
func (c *Client) Publish(ctx context.Context, topic string, message []byte) error {
	return nil
}

// Subscribe to do nothing.
func (c *Client) Subscribe(ctx context.Context, topic string, handlerFunc pubsub.HandlerFunc) error {
	return nil
}

// Close to do nothing.
func (c *Client) Close() error {
	return nil
}
