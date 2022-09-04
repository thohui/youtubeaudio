package mq

import (
	"github.com/streadway/amqp"
)

type Client struct {
	amqp      *amqp.Connection
	channel   *amqp.Channel
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
		channel:   channel,
		queueName: queueName,
	}, nil
}

func (c *Client) Consume() (<-chan amqp.Delivery, error) {
	q, err := c.channel.QueueDeclare(
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
	return c.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

func (c *Client) PublishResponse(replyTo, correlationID string, data []byte) {
	c.channel.Publish(
		"",      // exchange
		replyTo, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			CorrelationId: correlationID,
			ContentType:   "application/json",
			Body:          data,
		})

}
