package main

import (
	"io"
	"log"
	"log/slog"

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

func receive(ch *amqp.Channel, q amqp.Queue) {
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			slog.Info("Message received", "content", d.Body)
		}
	}()

	slog.Info("Waiting for messages, to exit press ^C")

	var forever chan struct{}
	<-forever
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

	receive(ch, q)
}
