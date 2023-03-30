package rabbitmq

import (
	"context"
	"email-serving-datathon/models"
	"encoding/json"
	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

func Connect() (*amqp.Connection, error) {
	return amqp.Dial(os.Getenv("AMQP_URI"))
}

func CreateChannel(connection *amqp.Connection) (*amqp.Channel, error) {
	return connection.Channel()
}

func InitAmqp() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(os.Getenv("AMQP_URI"))
	if err != nil {
		return &amqp.Connection{}, &amqp.Channel{}, nil
	}
	ch, err := conn.Channel()
	if err != nil {
		return &amqp.Connection{}, &amqp.Channel{}, nil
	}
	return conn, ch, nil

}

func CreateQueue(channel *amqp.Channel, name string) (amqp.Queue, error) {
	return channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func Publish(channel *amqp.Channel, queue amqp.Queue, body []byte) error {
	return channel.PublishWithContext(context.Background(),
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

func Consume(channel *amqp.Channel, queue amqp.Queue) (<-chan amqp.Delivery, error) {
	return channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
}

func CodeMessage(message models.Message) ([]byte, error) {
	return json.Marshal(message)
}

func DecodeMessage(body []byte) (models.Message, error) {
	var message models.Message
	err := json.Unmarshal(body, &message)
	if err != nil {
		return models.Message{}, err
	}

	return message, nil
}
