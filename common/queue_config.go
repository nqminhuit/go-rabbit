package common

import (
	"log/slog"
	"server/utils"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareQueue(ch *amqp.Channel) string {
	q, err := ch.QueueDeclare(
		"mdcore-reports",
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg:     amqp.QueueTypeQuorum,
			amqp.QueueMaxLenArg:   100_000,
			amqp.QueueOverflowArg: "reject-publish",
		})
	utils.FailOnError(err, "Failed to declare queue")
	return q.Name
}

func Connect(url string) *amqp.Connection {
	retry := 10
	for {
		conn, err := amqp.Dial(url)
		if err != nil {
			if retry < 1 {
				utils.FailOnError(err, "Failed to connect to RabbitMQ")
			}
			retry--
			time.Sleep(time.Second)
			continue
		} else {
			slog.Info("Connected to RabbitMQ")
			return conn
		}
	}
}
