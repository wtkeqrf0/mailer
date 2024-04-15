package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// Define the structure for the email.
type Parsable struct {
	To          []string
	Subject     string
	CopyTo      []string       // the recipient is explicitly know about copy.
	BlindCopyTo []string       // the recipient is explicitly don't know about copy.
	Sender      string         // set another email sender.
	ReplyTo     string         // to whom the recipient will respond.
	Parts       []Part         // message body parts.
	PartValues  map[string]any // used only with part body.
}

// Define the structure for the email parts.
type Part struct {
	ContentType int    `json:"content_type"`
	Body        []byte `json:"body"`
}

func main() {
	// Sample data
	email := Parsable{
		To:      []string{"artemka111777@mail.ru", "kirillsafatov@yandex.ru"},
		Subject: "hello",
		Parts: []Part{
			{ContentType: 0, Body: []byte("Ya sosal, menya ebali")},
		},
	}

	// Convert the email struct to JSON
	emailJSON, err := json.Marshal(email)
	if err != nil {
		log.Fatalf("Error marshalling JSON: %s", err)
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://admin:EKKZPGpx68KHkki@194.58.94.69:5678")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Declare a queue
	q, err := ch.QueueDeclare(
		"email", // Name of the queue
		true,    // Durable
		false,   // Delete when unused
		false,   // Exclusive
		false,   // No-wait
		amqp.Table{
			"x-queue-type":     "quorum", // For delivery limit
			"x-delivery-limit": 5,
		}, // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	// Publish the email JSON to the queue
	err = ch.Publish(
		"",     // Exchange
		q.Name, // Routing key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        emailJSON,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}

	log.Println("Message sent to the queue")
}
