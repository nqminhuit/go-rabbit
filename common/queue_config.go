package common

import (
	"server/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareQueue(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(
		"mdcore-reports",
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg:     amqp.QueueTypeQuorum,
			amqp.QueueMaxLenArg:   10,
			amqp.QueueOverflowArg: "reject-publish",
		})
	utils.FailOnError(err, "Failed to declare queue")
	return q
}
