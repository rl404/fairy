// Package kafka is a wrapper of the original "github.com/segmentio/kafka-go" library.
package kafka

import (
	"context"
	"time"

	"github.com/rl404/fairy/pubsub"
	"github.com/segmentio/kafka-go"
)

// Client is kafka pubsub client.
type Client struct {
	url         string
	middlewares []func(pubsub.HandlerFunc) pubsub.HandlerFunc
	writer      *kafka.Writer
}

// New to create new kafka pubsub client.
func New(url string) (*Client, error) {
	return &Client{
		url: url,
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(url),
			AllowAutoTopicCreation: true,
		},
	}, nil
}

// Use to add pubsub middlewares.
func (c *Client) Use(middlewares ...func(pubsub.HandlerFunc) pubsub.HandlerFunc) {
	c.middlewares = append(c.middlewares, middlewares...)
}

func (c *Client) applyMiddlewares(handlerFunc pubsub.HandlerFunc) pubsub.HandlerFunc {
	for _, mw := range c.middlewares {
		handlerFunc = mw(handlerFunc)
	}
	return handlerFunc
}

// Publish to publish message.
func (c *Client) Publish(ctx context.Context, topic string, data []byte) error {
	return c.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: data,
	})
}

// Subscribe to subscribe topic.
func (c *Client) Subscribe(ctx context.Context, topic string, handlerFunc pubsub.HandlerFunc) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{c.url},
		Topic:          topic,
		GroupBalancers: []kafka.GroupBalancer{kafka.RoundRobinGroupBalancer{}, kafka.RangeGroupBalancer{}},
	})

	go func(r *kafka.Reader, h pubsub.HandlerFunc) {
		h = c.applyMiddlewares(h)

		for {
			if err := reader.SetOffsetAt(ctx, time.Now()); err != nil {
				return
			}

			msg, err := r.ReadMessage(ctx)
			if err != nil {
				return
			}

			h(ctx, msg.Value)
		}
	}(reader, handlerFunc)

	return nil
}

// Close to close pubsub connection.
func (c *Client) Close() error {
	return c.writer.Close()
}
