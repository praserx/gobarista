// Copyright 2023 PraserX
package logger

import (
	"os"

	"log/slog"
)

// GlobalLogger global package instance.
var glog *slog.Logger

func init() {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	glog = slog.New(jsonHandler)
}

func Info(msg string) {
	glog.Info(msg)
}

func Warning(msg string) {
	glog.Warn(msg)
}

func Error(msg string) {
	glog.Error(msg)
}
