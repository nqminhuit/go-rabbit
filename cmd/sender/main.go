package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"server/service"
	"server/utils"
)

func main() {
	mq := &service.RabbitMQ{
		Url:            "amqp://guest:guest@localhost:5672",
		QueueName:      "mdcorereports",
		Exchange:       "",
		QueueMaxLenArg: 100_000,
	}
	mq.Connect()
	defer mq.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /mdcore/integration/console/{deploymentId}/report/scan", service.ReceiveReportHandler(mq))

	port := os.Getenv("SVPORT")
	slog.Info("Server is up and running")
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		utils.FailOnError(err, "Could not create http server")
	}
}
