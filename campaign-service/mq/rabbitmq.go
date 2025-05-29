package mq

import (
	"context"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
)

// refactor error handling
func handleError(err error, msg string){
	if err != nil{
		log.Fatalf("%s: %s", msg, err)
	}
}
func InitRabbitMQ(){
	var err error
	// Get the rabbitmq URL from env
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	// Setting up a connection to RabbitMQ
	conn, err = amqp.Dial(rabbitMQURL)
	handleError(err, "Can't connect to AMQP")

	// If connection succeed, establish a channel, channel serves as the communication protocol over the connection
	ch, err = conn.Channel()
	handleError(err, "Can't create amqpChannel")

	// Declare a new exchange that later will push messages to queues
	// For the exchange type, we use "direct" type to send messages to queues by the exact matching on the routing key
	err = ch.ExchangeDeclare(
		"campaign-events", //name
		"direct", // type
		true,     // durable
  		false,    // auto-deleted
  		false,    // internal
  		false,    // no-wait
  		nil,      // arguments
	)
	handleError(err,"Failed to declare an exhange")

}

func PublishCampaign(body []byte, routingKey string){
	// Create a context for connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// now we publish our exchange
	err := ch.PublishWithContext(ctx,
		"campaign-events", // exchange
		routingKey,     // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
	})
	handleError(err,"Failed to publish a message")
	log.Printf(" [x] sent %v", body)
}
