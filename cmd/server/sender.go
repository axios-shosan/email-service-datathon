package main

import (
	"email-serving-datathon/infra/rabbitmq"
	"email-serving-datathon/models"
	"email-serving-datathon/utils"
	_ "github.com/joho/godotenv/autoload"
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

	msg := models.Message{
		Email:    "mohamedzaouidi@yahoo.com",
		TeamCode: 1,
		UserId:   2,
		Counter:  0,
	}

	msgCoded, err := rabbitmq.CodeMessage(msg)

	utils.Panic(err, "Failed to encode message")

	err = rabbitmq.Publish(ch, confirmationQ, msgCoded)

	utils.Panic(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s\n", msgCoded)
}
