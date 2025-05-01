package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"server/common"
	"server/utils"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBatchConsumer struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Messages   []amqp.Delivery
	BatchSize  int
	QueueName  string
	LastProcess time.Time
	InactiveTimeoutSecond int
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
		go func() {
			for range time.Tick(time.Duration(consumer.InactiveTimeoutSecond) * time.Second) {
				msgSize := len(consumer.Messages)
				if sinceSecond := time.Since(consumer.LastProcess).Seconds();
					sinceSecond > float64(consumer.InactiveTimeoutSecond) && msgSize > 0 {
					slog.Info("consumming left over messages", "size", msgSize, "sinceSecond", sinceSecond)
					consumer.process()
				}
			}
		}()
		for d := range msgs {
			consumer.Messages = append(consumer.Messages, d)
			if len(consumer.Messages) >= consumer.BatchSize {
				slog.Info("consumming batch", "size", len(consumer.Messages))
				consumer.process()
			}
		}
	}()

	slog.Info("Waiting for messages, to exit press ^C")

	forever := make(chan struct{})
	<-forever
}

func (consumer *RabbitMQBatchConsumer) process() {
	slog.Info("processing")
	i := 0
	var lastMsg *amqp.Delivery
	for _, d := range consumer.Messages {
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
		// make a buffer channel with size = 50
		// send json to that buffer channel
		// ack

		i++
		lastMsg = &d
	}

	err := lastMsg.Ack(true)
	utils.LogOnError(err, "Could not ack message")

	consumer.Messages = []amqp.Delivery{}
	consumer.LastProcess = time.Now()
	slog.Info("lastProcessed", "time", consumer.LastProcess)
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer utils.Close(conn)

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a RabbitMQ channel")
	defer utils.Close(ch)

	qName := common.DeclareQueue(ch)

	consumer := &RabbitMQBatchConsumer{
		Connection: conn,
		BatchSize:  5,
		QueueName:  qName,
		Channel:    ch,
		InactiveTimeoutSecond: 3,
	}
	consumer.start()
}
