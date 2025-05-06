package http

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"server/transport/amqp"
	"server/utils"
	"time"

	proto "server/transport/grpc/proto"
)

func ReceiveReportHandler(mq *amqp.RabbitMQ, grpcClient *proto.GrpcClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		dataRetentionInMillis, err := grpcClient.GetDataRetentionMillis(ctx, "afe43c63-1d83-4b55-84c8-2a71dea3ea41")
		if err != nil {
			slog.Error("Failed to get data from account grpc service", "Reason", err.Error())
			return
		}

		w.Header().Add("dataRetentionInMillis", dataRetentionInMillis)
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
