// Package rabbitmq is a wrapper of the original "github.com/streadway/amqp" library.
//
// todo: add reconnect feature.
package rabbitmq

import (
	"context"

	pubsub "github.com/rl404/fairy/pubsub_v2"
	"github.com/streadway/amqp"
)

// Client is rabbitmq pubsub client.
type Client struct {
	client *amqp.Connection
}

// New to create new rabbitmq pubsub client.
func New(url string) (*Client, error) {
	c, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Client{client: c}, nil
}

// Publish to publish message.
func (c *Client) Publish(ctx context.Context, queue string, data []byte) error {
	ch, err := c.client.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	if err := ch.Publish("", q.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	}); err != nil {
		return err
	}

	return nil
}

// Subscribe to subscribe queue.
func (c *Client) Subscribe(ctx context.Context, queue string, handlerFunc pubsub.HandlerFunc) error {
	ch, err := c.client.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func(msgs <-chan amqp.Delivery) {
		for msg := range msgs {
			handlerFunc(ctx, msg.Body)
		}
	}(msgs)

	return nil
}

// Close to close pubsub connection.
func (c *Client) Close() error {
	return c.client.Close()
}
