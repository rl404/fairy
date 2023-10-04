// Package redis is a wrapper of the original "github.com/redis/go-redis" library.
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rl404/fairy/pubsub"
)

// Client is redis pubsub client.
type Client struct {
	client *redis.Client
}

// New to create new redis pubsub client.
func New(address, password string) (*Client, error) {
	return NewWithConfig(redis.Options{
		Addr:     address,
		Password: password,
	})
}

// NewWithConfig to create pubsub from go-redis options.
func NewWithConfig(option redis.Options) (*Client, error) {
	client := redis.NewClient(&option)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping test.
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return NewFromGoRedis(client), nil
}

// NewFromGoRedis to create pubsub from go-redis client.
func NewFromGoRedis(client *redis.Client) *Client {
	return &Client{
		client: client,
	}
}

// Publish to publish message.
func (c *Client) Publish(ctx context.Context, channel string, data []byte) error {
	return c.client.Publish(ctx, channel, data).Err()
}

// Subscribe to subscribe channel.
func (c *Client) Subscribe(ctx context.Context, channel string, handlerFunc pubsub.HandlerFunc) error {
	ch := c.client.Subscribe(ctx, channel)

	go func(c *redis.PubSub) {
		for msg := range c.Channel() {
			handlerFunc(ctx, []byte(msg.Payload))
		}
	}(ch)

	return nil
}

// Close to close redis pubsub client.
func (c *Client) Close() error {
	return c.client.Close()
}
