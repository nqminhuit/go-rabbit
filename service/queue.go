package service

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"mime/multipart"
	"server/utils"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Url            string
	Conn           *amqp.Connection
	QueueName      string
	QueueMaxLenArg int
	Exchange       string
}

func (mq *RabbitMQ) EnsureQueue() *amqp.Channel {
	ch, err := mq.Conn.Channel()
	utils.FailOnError(err, "Failed to open a RabbitMQ channel")

	_, err = ch.QueueDeclare(
		mq.QueueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			amqp.QueueTypeArg:     amqp.QueueTypeQuorum,
			amqp.QueueMaxLenArg:   mq.QueueMaxLenArg,
			amqp.QueueOverflowArg: "reject-publish",
		})
	utils.FailOnError(err, "Failed to declare queue")
	return ch
}

func (mq *RabbitMQ) Send(part *multipart.Part) error {
	defer utils.Close(part)

	ch := mq.EnsureQueue() // run on every api request???
	defer utils.Close(ch)

	ack, nack := ch.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	err := ch.Confirm(false)
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
		mq.Exchange,
		mq.QueueName,
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

func (mq *RabbitMQ) Connect() {
	retry := 60
	for {
		conn, err := amqp.Dial(mq.Url)
		if err != nil {
			if retry < 1 {
				utils.FailOnError(err, "Failed to connect to RabbitMQ")
			}
			retry--
			time.Sleep(time.Second)
			continue
		} else {
			slog.Info("Connected to RabbitMQ")
			mq.Conn = conn
			return
		}
	}
}

func (mq *RabbitMQ) Close() {
	utils.Close(mq.Conn)
}
