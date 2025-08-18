package main

import (
	"fmt"
	"local-chat/internal/message"
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
