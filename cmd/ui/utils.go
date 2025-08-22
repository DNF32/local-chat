package main

import (
	"fmt"
	"local-chat/internal/message"
	"strings"
	"time"
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

func MakeMsgFromType(username string, msgType message.MessageType, content string) message.Message {
	var msg message.Message
	switch msgType {
	case message.Text:
		msg = message.Message{Type: message.Text,
			Username:  username,
			Content:   content,
			Timestamp: time.Now()}
	case message.Join:
		msg = message.Message{Type: message.Join,
			Username:  username,
			Content:   content,
			Timestamp: time.Now()}
	case message.Leave:
		msg = message.Message{Type: message.Leave,
			Username:  username,
			Content:   content,
			Timestamp: time.Now()}
	}
	return msg
}

