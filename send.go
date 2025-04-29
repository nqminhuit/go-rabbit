package main

import (
	"context"
	"fmt"
	"log/slog"
	"server/common"
	"server/utils"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func send(ch *amqp.Channel, q amqp.Queue, msg string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(msg),
		})
	utils.LogOnError(err, "Could not publilsh a message")
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer utils.Close(conn)

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a RabbitMQ channel")
	defer utils.Close(ch)

	q := common.DeclareQueue(ch)

	ack, nack := ch.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	i := 0
	for range time.Tick(10 * time.Millisecond) {
		msg := fmt.Sprintf("hello batminh %d", i)
		err = ch.Confirm(false)
		if err != nil {
			slog.Error("Could not confirm", "Reason", err.Error())
		}
		send(ch, q, msg)

		select {
		case <-ack:
			slog.Info("Message sent", "content", msg)
		case <-nack:
			slog.Error("Consumers overload, slowing down")
			time.Sleep(3 * time.Second)
		}
		i++
	}
}
