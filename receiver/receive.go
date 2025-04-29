package main

import (
	"io"
	"log"
	"log/slog"
	"math/rand/v2"
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

func receive(ch *amqp.Channel, q amqp.Queue) {
	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			processingTimeMs := rand.IntN(500)
			time.Sleep(time.Duration(processingTimeMs) * time.Millisecond)
			slog.Info("Message processed", "processingTimeMs", processingTimeMs, "content", d.Body)
			d.Ack(false)
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
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg:     amqp.QueueTypeQuorum,
			amqp.QueueMaxLenArg:   10,
			amqp.QueueOverflowArg: "reject-publish",
		})
	failOnError(err, "Failed to declare queue")

	err = ch.Qos(1, 0, false)

	receive(ch, q)
}
