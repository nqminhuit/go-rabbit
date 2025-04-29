package main

import (
	"context"
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
			ContentType: "text/plain",
			Body:        []byte(msg)})
	logOnError(err, "Could not publilsh a message")
	slog.Info("Message sent", "content", msg)
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
		false, false,
		nil)
	failOnError(err, "Failed to declare queue")

	send(ch, q, "hello batminh!!")
}
