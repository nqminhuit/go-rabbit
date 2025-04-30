package main

import (
	"context"
	"log/slog"
	"net/http"
	"server/common"
	"server/utils"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func send(ch *amqp.Channel, q amqp.Queue, msg string) {
	err := ch.Confirm(false)
	if err != nil {
		slog.Error("Could not confirm", "Reason", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
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

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /mdcore/integration/console/{deploymentId}/report/scan",
		func(w http.ResponseWriter, r *http.Request) {
			msg := r.PathValue("deploymentId")
			send(ch, q, msg)
			select {
			case <-ack:
				w.WriteHeader(200)
			case <-nack:
				w.WriteHeader(429)
			}
		})

	slog.Info("Server is up and listening on port 9093")
	err = http.ListenAndServe(":9093", mux)
	if err != nil {
		utils.FailOnError(err, "Could not create http server")
	}
}
