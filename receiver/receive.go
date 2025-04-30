package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"server/common"
	"server/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func compactJson(msgs <-chan amqp.Delivery) <-chan *bytes.Buffer {
	out := make(chan *bytes.Buffer)
	go func() {
		for msg := range msgs {
			compacted := &bytes.Buffer{}
			err := json.Compact(compacted, msg.Body)
			utils.LogOnError(err, "Could not process json")

			err = msg.Ack(false)
			utils.LogOnError(err, "Could not ack message")
			out <- compacted
		}
		close(out)
	}()
	return out
}

func addMetadata(msgs <-chan *bytes.Buffer) <-chan map[string]any {
	out := make(chan map[string]any)
	go func() {
		var data map[string]any
		for msg := range msgs {
			err := json.Unmarshal(msg.Bytes(), &data)
			utils.LogOnError(err, "Could not unmarshal json")
			out <- data
		}
		close(out)
	}()
	return out
}

func batch(data <-chan map[string]any, batchSize int) <-chan []map[string]any {
	out := make(chan []map[string]any)
	go func() {
		batched := make([]map[string]any, 0, batchSize)
		for msg := range data {
			if len(batched) < batchSize {
				batched = append(batched, msg)
			} else {
				out <- batched
				slog.Info("batch", "batched", batched)
				batched = append(make([]map[string]any, 0, batchSize), msg)
			}
		}
		out <- batched
		close(out)
	}()
	return out
}

func sendDataToOpenSearch(data <-chan []map[string]any) {
	go func() {
		for msg := range data {
			_ = msg
			// slog.Info("sendDataToOpenSearch", "msg", msg)
		}
	}()
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
	utils.FailOnError(err, "Failed to register a consumer")

	compactJson := compactJson(msgs)
	addedMetadata := addMetadata(compactJson)
	batched := batch(addedMetadata, 5)
	sendDataToOpenSearch(batched)

	slog.Info("Waiting for messages, to exit press ^C")

	forever := make(chan struct{})
	<-forever
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer utils.Close(conn)

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a RabbitMQ channel")
	defer utils.Close(ch)

	q := common.DeclareQueue(ch)

	err = ch.Qos(20, 0, false)
	utils.LogOnError(err, "Failed to config fair dispatch on channel")

	receive(ch, q)
}
