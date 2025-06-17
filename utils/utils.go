package utils

import (
	"fmt"
	"io"
	"log"
	"log/slog"
)

func FailOnError(err error, msg string, args ...string) {
	if err != nil {
		log.Panicf("%s: %s", fmt.Sprintf(msg, args), err)
	}
}

func LogOnError(err error, msg string, args ...string) {
	if err != nil {
		slog.Error(fmt.Sprintf(msg, args), "Reason", err.Error())
	}
}

func Close(c io.Closer) {
	err := c.Close()
	LogOnError(err, "Could not close")
}
