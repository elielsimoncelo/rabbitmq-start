package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func bodyGenerator() string {
	return fmt.Sprintf("Logging message at %s", time.Now().Format("2006-01-02 15:04:05"))
}

func main() {
	connection, err := amqp.Dial("amqp://rabbitmq:rabbitmq@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer connection.Close()

	channel, err := connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	err = channel.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	for {
		body := bodyGenerator()

		err = channel.Publish(
			"logs",     // exchange
			"log.file", // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			},
		)

		failOnError(err, "Failed to publish a message")

		log.Printf(" [x] Sent %s", body)

		time.Sleep(time.Second * 1)
	}
}
