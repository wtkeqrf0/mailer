package rabbit

import (
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

// Connection represents a RabbitMQ connection instance
type Connection struct {
	channel *amqp.Channel
}

// NewConn creates a new RabbitMQ producer instance
func NewConn(ctx context.Context, url string) *Connection {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	context.AfterFunc(ctx, func() {
		if err = conn.Close(); err != nil {
			log.Println(err.Error())
		}
	})

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &Connection{
		channel: channel,
	}
}

type Produce func(msg json.RawMessage) error

// Publisher sends a message to a specified exchange with a routing key
func (r *Connection) Publisher(queueName string) Produce {
	if queue, err := r.channel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,
	); err != nil {
		panic(err)
	} else {
		queueName = queue.Name
	}

	return func(msg json.RawMessage) error {
		return r.channel.Publish(
			"",        // Exchange
			queueName, // Routing key
			false,     // Mandatory
			false,     // Immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        msg,
			},
		)
	}
}

// Consumer consumes messages from a specified exchange with a routing key
func (r *Connection) Consumer(ctx context.Context, queueName string) <-chan amqp.Delivery {
	queue, err := r.channel.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		amqp.Table{
			"x-queue-type":     "quorum", // For delivery limit
			"x-delivery-limit": 5,
		}, // Arguments
	)
	if err != nil {
		panic(err)
	}

	const consumerName = `main-consumer`
	ch, err := r.channel.Consume(
		queue.Name,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	context.AfterFunc(ctx, func() {
		if err = r.channel.Cancel(consumerName, false); err != nil {
			log.Println(err.Error())
		}
	})
	return ch
}
