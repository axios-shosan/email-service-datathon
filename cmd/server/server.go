package main

import (
	database "email-serving-datathon/infra/db"
	"email-serving-datathon/infra/mailer"
	"email-serving-datathon/infra/rabbitmq"
	"email-serving-datathon/utils"
	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {

	auth := mailer.AuthSmtp()
	db := database.Connect()
	database.Migrate(db)

	conn, ch, err := rabbitmq.InitAmqp()
	utils.Panic(err, "Failed to Init Amqp")
	defer func(conn *amqp.Connection) {
		err = conn.Close()
		utils.Panic(err, "Failed to close connection to RabbitMQ")
	}(conn)
	defer func(ch *amqp.Channel) {
		err := ch.Close()
		utils.Panic(err, "Failed to close rabbit mq channel")
	}(ch)

	confirmationQ, err := rabbitmq.CreateQueue(ch, "confirmation-queue")
	utils.Panic(err, "Failed to declare queue")

	msgs, err := rabbitmq.Consume(ch, confirmationQ)

	utils.Panic(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		mailer.TreatMessages(auth, ch, confirmationQ, db, msgs)
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
