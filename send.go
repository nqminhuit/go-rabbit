package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"server/common"
	"server/utils"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func send(conn *amqp.Connection, part *multipart.Part) error {
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a RabbitMQ channel")
	defer utils.Close(ch)

	qName := common.DeclareQueue(ch) // declare queue every api request???

	ack, nack := ch.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	defer utils.Close(part)
	err = ch.Confirm(false)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	buffer := new(bytes.Buffer)
	_, err = buffer.ReadFrom(part)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(
		ctx,
		"",
		qName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         buffer.Bytes(),
		})
	if err != nil {
		return err
	}

	select {
	case <-ack:
		return nil
	case <-nack:
		return errors.New("429")
	}
}

func processSingleReport(reader *multipart.Reader) (*multipart.Part, error) {
	part, err := reader.NextPart()
	if err != nil {
		return nil, err
	}
	return part, nil
}

func main() {
	conn := common.Connect("amqp://guest:guest@localhost:5672")
	defer utils.Close(conn)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /mdcore/integration/console/{deploymentId}/report/scan",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("dataRetentionInMillis", "1234567890")
			// depId := r.PathValue("deploymentId")
			reader, err := r.MultipartReader()
			utils.LogOnError(err, "Could not get multipart reader")

			for {
				part, err := processSingleReport(reader)
				if err == io.EOF {
					break
				}
				if err != nil {
					slog.Error("Could not read part", "Reason", err.Error())
					continue
				}

				err = send(conn, part)
				utils.LogOnError(err, "Could not send message to RabbitMQ")
				if err != nil && err.Error() == "429" {
					w.WriteHeader(429)
					return
				}
			}
			w.WriteHeader(200)
		})

	slog.Info("Server is up and listening on port 9093")
	err := http.ListenAndServe(":9093", mux)
	if err != nil {
		utils.FailOnError(err, "Could not create http server")
	}
}
