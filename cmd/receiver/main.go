package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"server/service"
	mq "server/transport/amqp"
	"server/utils"
	"strings"

	"github.com/opensearch-project/opensearch-go/v4/opensearchutil"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBatchConsumer struct {
	Connection       *amqp.Connection
	Channel          *amqp.Channel
	BatchSize        int
	QueueName        string
	OpenSearchClient *service.OpenSearchClient
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

	var lastMsg *amqp.Delivery
	totalItems := new(int)

	opensearch := consumer.OpenSearchClient
	indexer, err := opensearchutil.NewBulkIndexer(opensearchutil.BulkIndexerConfig{
		Client:     opensearch.Client,
		NumWorkers: 1,
		FlushBytes: 1e+7,
		Index:      opensearch.IndexName,
		Pipeline:   service.INGEST_PIPELINE_NAME,
		OnFlushEnd: func(_ context.Context) {
			slog.Info("Flushed", "totalItems", *totalItems)
			*totalItems = 0
			if lastMsg != nil {
				err = lastMsg.Ack(true)
				utils.LogOnError(err, "Could not ack message")
			}
		},
	})
	utils.LogOnError(err, "Failed to create bulk indexer")

	go func() {
		for d := range msgs {
			// 1. add metadata to json: TODO
			var data map[string]any
			err = json.Unmarshal(d.Body, &data)
			utils.LogOnError(err, "Could not unmarshal json")
			dataId := data["data_id"].(string)
			data["deploymentId"] = "5ccab8bb-ceaa-41e1-bd56-fb9660978843"
			data["accountId"] = "04a4b948-6b80-4eee-bb99-c4495e7db415"
			data["instanceId"] = "85ddeecd-7448-4cb8-b8af-f6a0d5fcff1b"
			data["instanceName"] = "ceafdb4d-43e2-4a29-9f41-ed00dc7b7d14"
			data["groupId"] = "2cb6ea22-1f58-4002-b961-e3e0cb485fb2"
			data["groupName"] = "01df93b1-ccf7-4627-b008-0d0aea1531f8"
			data["retentionMs"] = "2c739207-2b44-42fc-b5b7-eef2b58b8e5f"

			// 2. marshal to json
			content, err := json.Marshal(&data)
			utils.LogOnError(err, "Failed to marshal data")

			// 3. send json to opensearch
			consumer.OpenSearchClient.AddToBulk(&indexer, dataId, bytes.NewReader(content))
			*totalItems++

			lastMsg = &d
		}
	}()

	slog.Info("Waiting for messages...")

	forever := make(chan struct{})
	<-forever
}

func main() {
	mq := &mq.RabbitMQ{
		Url:            "amqp://guest:guest@localhost:5672",
		QueueName:      "mdcorereports",
		Exchange:       "",
		QueueMaxLenArg: 100_000,
	}
	mq.Connect()
	defer mq.Close()

	ch := mq.EnsureQueue()
	defer utils.Close(ch)

	coreIndexName := os.Getenv("OPENSEARCH_INDEX_NAME_MDCORE")
	username := os.Getenv("OPENSEARCH_USERNAME")
	password := os.Getenv("OPENSEARCH_PASSWORD")
	addresses := strings.Split(os.Getenv("OPENSEARCH_ADDRESSES"), ",")

	o := service.ConnectToOpenSearch(coreIndexName, username, password, addresses)

	consumer := &RabbitMQBatchConsumer{
		Connection: mq.Conn,
		BatchSize:  1000,
		QueueName:  mq.QueueName,
		Channel:    ch,
		OpenSearchClient: &service.OpenSearchClient{
			Client:    o.Client,
			IndexName: coreIndexName,
		},
	}
	consumer.start()
}
