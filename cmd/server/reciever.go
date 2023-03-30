package main

import (
	"email-serving-datathon/infra/rabbitmq"
	"email-serving-datathon/utils"
	"log"
)

func main() {
	conn, err := rabbitmq.Connect()
	utils.Panic(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := rabbitmq.CreateChannel(conn)
	utils.Panic(err, "Failed to open a channel")
	defer ch.Close()

	confirmationQ, err := rabbitmq.CreateQueue(ch, "confirmation-queue")
	utils.Panic(err, "Failed to declare queue")

	msgs, err := rabbitmq.Consume(ch, confirmationQ)

	utils.Panic(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			msg, err := rabbitmq.DecodeMessage(d.Body)
			utils.Panic(err, "Failed to decode the message")
			log.Printf("Participant Email: %s", msg.Email)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
