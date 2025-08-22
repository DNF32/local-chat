package logger

import (
	"log/slog"
	"os"
)

const UI_LOG_PATH = "/Users/dnf/code/local-chat/ui_debug.log"
const SERVER_LOG_PATH = "/Users/dnf/code/local-chat/backend_debug.log"

func NewFileLogger(path string) (*slog.Logger, error) {
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		return nil, err
	}

	handler := slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.Level(-4)})
	logger := slog.New(handler)
	return logger, nil
}
