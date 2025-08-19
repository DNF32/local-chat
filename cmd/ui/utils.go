package main

import (
	"fmt"
	"local-chat/internal/message"
	"log/slog"
	"os"
	"strings"
)

func ParseMsgType(s string) (message.MessageType, error) {

	switch s[0] {
	case '/':
		if strings.HasPrefix(s, "/join") {
			return message.Join, nil
		} else if strings.HasPrefix(s, "/leave") {
			return message.Leave, nil
		} else {
			return message.Undefined, fmt.Errorf("Invalid command sent")
		}
	default:
		return message.Text, nil
	}
}

func NewFileLogger(path string) (*slog.Logger, error) {
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		return nil, err
	}

	handler := slog.NewTextHandler(file, nil)
	logger := slog.New(handler)
	return logger, nil
}
