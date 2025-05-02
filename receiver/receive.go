package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"server/common"
	"server/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBatchConsumer struct {
	Connection            *amqp.Connection
	Channel               *amqp.Channel
	BatchSize             int
	QueueName             string
}

func (consumer *RabbitMQBatchConsumer) start() {
	ch := consumer.Channel

	err := ch.Qos(consumer.BatchSize, 0, false)
	utils.LogOnError(err, "Failed to config fair dispatch on channel")

	msgs, err := ch.Consume(
		consumer.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil)
	utils.FailOnError(err, "Failed to register a consumer")

	go func() {
		i := 0
		var lastMsg *amqp.Delivery
		for d := range msgs {
			// 1. compact json
			compacted := &bytes.Buffer{}
			err := json.Compact(compacted, d.Body)
			utils.LogOnError(err, "Could not process json")

			var data map[string]any
			err = json.Unmarshal(compacted.Bytes(), &data)
			utils.LogOnError(err, "Could not unmarshal json")
			slog.Info("compacted", "dataId", data["data_id"], "i", i)

			// 2. add metadata to json:
			// "deploymentId", deploymentId,
			// "accountId", accountId,
			// "instanceId", instanceId,
			// "instanceName", SecurityContextUtils.getInstanceName(),
			// "groupId", groupId,
			// "groupName", SecurityContextUtils.getGroupName(),
			// "retentionMs", retentionMs

			// 3. send to opensearch

			i++
			lastMsg = &d
			if i == consumer.BatchSize {
				err := lastMsg.Ack(true)
				utils.LogOnError(err, "Could not ack message")
				i = 0
			}
		}
	}()

	slog.Info("Waiting for messages...")

	forever := make(chan struct{})
	<-forever
}

func main() {
	conn := common.Connect("amqp://guest:guest@localhost:5672/")
	defer utils.Close(conn)

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a RabbitMQ channel")
	defer utils.Close(ch)

	qName := common.DeclareQueue(ch)

	consumer := &RabbitMQBatchConsumer{
		Connection:            conn,
		BatchSize:             50,
		QueueName:             qName,
		Channel:               ch,
	}
	consumer.start()
}
