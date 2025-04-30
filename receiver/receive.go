package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"server/common"
	"server/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func sendDataToOpenSearch() {
	// listen to a buffer channel
	// when the channel reach max capacity
	// send data to opensearch
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

	go func() {
		i := 0
		for d := range msgs {
			// 1. compact json
			compacted := &bytes.Buffer{}
			err := json.Compact(compacted, d.Body)
			utils.LogOnError(err, "Could not process json")

			// 2. add metadata to json:
			// "deploymentId", deploymentId,
			// "accountId", accountId,
			// "instanceId", instanceId,
			// "instanceName", SecurityContextUtils.getInstanceName(),
			// "groupId", groupId,
			// "groupName", SecurityContextUtils.getGroupName(),
			// "retentionMs", retentionMs

			// 3. send to opensearch
			// make a buffer channel with size = 50
			// send json to that buffer channel
			// ack

			slog.Info("Message processed", "i", i)
			err = d.Ack(false)
			utils.LogOnError(err, "Could not ack message")
			i++
		}
	}()

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
