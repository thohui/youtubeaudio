package mq

import (
	"encoding/json"
	"sync"

	"github.com/lithammer/shortuuid"
	"github.com/streadway/amqp"
	"github.com/thohui/youtubeaudio/internal/structures"
)

type Client struct {
	amqp         *amqp.Connection
	channel      *amqp.Channel
	queueName    string
	rpcQueueName string
	jobs         map[string]chan<- []byte
	mutex        sync.Mutex
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
		amqp:         conn,
		channel:      channel,
		queueName:    queueName,
		rpcQueueName: "rpc_" + shortuuid.New(),
		jobs:         make(map[string]chan<- []byte),
		mutex:        sync.Mutex{},
	}, nil
}

func (c *Client) Publish(pl structures.BackendPublishPayload, job chan<- []byte) error {
	data, err := json.Marshal(pl)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	random := shortuuid.New()
	c.jobs[random] = job
	return c.channel.Publish(
		"",          // exchange
		c.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ReplyTo:       c.rpcQueueName,
			CorrelationId: random,
			ContentType:   "application/json",
			Body:          data,
		})
}

func (c *Client) HandleMessages() {
	q, err := c.channel.QueueDeclare(
		c.rpcQueueName, // name
		false,          // durable
		false,          // delete when unused
		true,           // exclusive
		false,          // noWait
		nil,            // arguments
	)
	if err != nil {
		panic(err)
	}
	messages, err := c.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}
	for msg := range messages {
		c.mutex.Lock()
		job, ok := c.jobs[msg.CorrelationId]
		if !ok {
			continue
		}
		job <- msg.Body
		close(job)
		delete(c.jobs, msg.CorrelationId)
		c.mutex.Unlock()
	}
}
