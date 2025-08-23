package logger

import (
	"fmt"
	"log/slog"
	"os"
)

var RepoPath string = ""

func init() {
	RepoPath = fmt.Sprintf("%s/%s", os.Getenv("HOME"), "code/local-chat/")
}

var UI_LOG_PATH = fmt.Sprintf(RepoPath, "ui_debug.log")
var SERVER_LOG_PATH = fmt.Sprintf(RepoPath, "backend_debug.log")

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
