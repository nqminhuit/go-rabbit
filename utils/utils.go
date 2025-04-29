package utils

import (
	"io"
	"log"
	"log/slog"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func LogOnError(err error, msg string) {
	if err != nil {
		slog.Error(msg, "Reason", err.Error())
	}
}

func Close(c io.Closer) {
	err := c.Close()
	LogOnError(err, "Could not close")
}
