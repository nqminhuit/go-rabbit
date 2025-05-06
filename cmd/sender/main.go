package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"server/transport/amqp"
	handler "server/transport/http"
	"server/utils"

	proto "server/transport/grpc/proto"
)

func main() {
	var err error
	mq := &amqp.RabbitMQ{
		Url:            "amqp://guest:guest@localhost:5672",
		QueueName:      "mdcorereports",
		Exchange:       "",
		QueueMaxLenArg: 100_000,
	}
	mq.Connect()
	defer mq.Close()

	grpc := &proto.GrpcClient{
		Url: ":50051",
	}
	grpc.Connect()
	defer grpc.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /mdcore/integration/console/{deploymentId}/report/scan", handler.ReceiveReportHandler(mq, grpc))

	port := os.Getenv("SVPORT")
	slog.Info("Server is up and running")
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		utils.FailOnError(err, "Could not create http server")
	}
}
