package mq

import (
	"github.com/streadway/amqp"
)

type Client struct {
	amqp      *amqp.Connection
	Channel   *amqp.Channel
	queueName string
}

func New(uri string, queueName string) (*Client, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Client{
		amqp:      conn,
		Channel:   channel,
		queueName: queueName,
	}, nil
}

func (c *Client) Close() error {
	return c.amqp.Close()
}

func (c *Client) Consume() (<-chan amqp.Delivery, error) {
	q, err := c.Channel.QueueDeclare(
		c.queueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return nil, err
	}
	return c.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}
