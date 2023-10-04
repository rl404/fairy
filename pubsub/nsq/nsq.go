// Package nsq is a wrapper of the original "github.com/nsqio/go-nsq" library.
package nsq

import (
	"context"

	"github.com/nsqio/go-nsq"
	"github.com/rl404/fairy/pubsub"
)

// Client is NSQ pubsub client.
type Client struct {
	address  string
	config   *nsq.Config
	producer *nsq.Producer
	consumer *nsq.Consumer
}

// New to create new NSQ pubsub client.
func New(address string) (*Client, error) {
	return NewWithConfig(address, nsq.NewConfig())
}

// NewWithConfig to create new NSQ pubsub client with config.
func NewWithConfig(address string, cfg *nsq.Config) (*Client, error) {
	producer, err := nsq.NewProducer(address, cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		address:  address,
		config:   cfg,
		producer: producer,
	}, nil
}

// Publish to publish message.
func (c *Client) Publish(ctx context.Context, topic string, data []byte) error {
	return c.producer.Publish(topic, data)
}

// Subscribe to subscriber topic.
func (c *Client) Subscribe(ctx context.Context, topic string, handlerFunc pubsub.HandlerFunc) (err error) {
	c.consumer, err = nsq.NewConsumer(topic, "channel", c.config)
	if err != nil {
		return err
	}

	c.consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		handlerFunc(context.Background(), msg.Body)
		return nil
	}))

	return c.consumer.ConnectToNSQD(c.address)
}

// Close to close pubsub connection.
func (c *Client) Close() error {
	if c.consumer != nil {
		c.consumer.Stop()
	}
	if c.producer != nil {
		c.producer.Stop()
	}
	return nil
}
