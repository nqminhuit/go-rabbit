package http

import (
	"io"
	"log/slog"
	"net/http"
	"server/transport/amqp"
	"server/utils"
)

func ReceiveReportHandler(mq *amqp.RabbitMQ) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("dataRetentionInMillis", "1234567890")
		reader, err := r.MultipartReader()
		if err != nil {
			slog.Error("Could not get multipart reader", "Reason", err.Error())
			return
		}

		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				slog.Error("Could not read part", "Reason", err.Error())
				continue
			}

			err = mq.Send(part)
			utils.LogOnError(err, "Could not send message to RabbitMQ")
			if err != nil && err.Error() == "429" {
				w.WriteHeader(429)
				return
			}
		}
		w.WriteHeader(200)
	}
}
