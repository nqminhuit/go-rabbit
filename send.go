package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func logOnError(err error, msg string) {
	if err != nil {
		slog.Error(msg, "Reason", err.Error())
	}
}

func close(c io.Closer) {
	err := c.Close()
	logOnError(err, "Could not close")
}

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
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	logOnError(err, "Could not publilsh a message")
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer close(conn)

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a RabbitMQ channel")
	defer close(ch)

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
	failOnError(err, "Failed to declare queue")

	ack, nack := ch.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	i := 0
	for range time.Tick(10 * time.Millisecond) {
		msg :=  fmt.Sprintf("hello batminh %d", i)
		err = ch.Confirm(false)
		if err != nil {
			slog.Error("Could not confirm", "Reason", err.Error())
		}
		send(ch, q,msg)

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
