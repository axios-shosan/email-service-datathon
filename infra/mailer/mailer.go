package mailer

import (
	"email-serving-datathon/infra/rabbitmq"
	"email-serving-datathon/models"
	"email-serving-datathon/services"
	"email-serving-datathon/utils"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"net/smtp"
	"os"
	"strconv"
)

func createReceiver(to string) []string {
	return []string{
		to,
	}
}
func createMsg(from, to, subject, content string) []byte {
	return []byte("From:" + from + "\r\n" +
		"To:" + to + "\r\n" +
		"Subject:" + subject + "\r\n\r\n" +
		content + "\r\n")
}

func SendMail(auth smtp.Auth, to, subject, content string) error {
	addr := os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT")
	from := os.Getenv("SMPT_SENDER_EMAIL")

	receiver := createReceiver(to)
	msg := createMsg("cse@esi.dz", to, subject, content)

	return smtp.SendMail(addr, auth, from, receiver, msg)
}

func SendConfirmation(auth smtp.Auth, to string, teamCode uint) error {
	fmt.Print()
	return SendMail(auth, to, "Cse Datathon Confirmation", "<h1>hello sir, here is your teamCode: "+strconv.FormatInt(int64(teamCode), 10))
}

func AuthSmtp() smtp.Auth {
	host := os.Getenv("SMTP_HOST")
	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	return smtp.PlainAuth("", user, password, host)
}

func TreatMessages(auth smtp.Auth, ch *amqp.Channel, queue amqp.Queue, db *gorm.DB, msgs <-chan amqp.Delivery) {
	for d := range msgs {
		go func(d amqp.Delivery) {
			msg, err := rabbitmq.DecodeMessage(d.Body)
			if err != nil {
				utils.Panic(err, "Failed to decode the message")
				err = services.SaveEmail(db, models.Mail{
					Model:    gorm.Model{},
					UserId:   msg.UserId,
					TeamCode: msg.TeamCode,
					Email:    msg.Email,
					Status:   models.Failed,
				})
				utils.Panic(err, "Failed to register email in database")
				return
			}

			err = SendConfirmation(auth, msg.Email, msg.TeamCode)
			if err != nil {
				err = HandleEmailFailure(ch, queue, db, msg)
				utils.Panic(err, "Error while handling Email Failure")
			}

			err = services.SaveEmail(db, models.Mail{
				Model:    gorm.Model{},
				UserId:   msg.UserId,
				TeamCode: msg.TeamCode,
				Email:    msg.Email,
				Status:   models.Success,
			})
			utils.Panic(err, "Failed to register email in database")

			fmt.Println("Email sent successfully")
		}(d)

	}
}

func HandleEmailFailure(ch *amqp.Channel, queue amqp.Queue, db *gorm.DB, msg models.Message) error {

	if msg.Counter < 3 {
		msgCoded, err := rabbitmq.CodeMessage(models.Message{
			Email:    msg.Email,
			TeamCode: msg.TeamCode,
			UserId:   msg.UserId,
			Counter:  msg.Counter + 1,
		})
		if err != nil {
			return err
		}
		return rabbitmq.Publish(ch, queue, msgCoded)
	} else {
		// admit the failure and just store it as it is
		return services.SaveEmail(db, models.Mail{
			Model:    gorm.Model{},
			UserId:   msg.UserId,
			TeamCode: msg.TeamCode,
			Email:    msg.Email,
			Status:   models.Failed,
		})
	}
}
