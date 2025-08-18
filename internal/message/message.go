package message

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type MessageType string

const (
	InitUser  MessageType = "initUser"
	Join      MessageType = "joined" // using this because we take advantage of the string method
	Text      MessageType = "send"
	Leave     MessageType = "left" // using this because we take advantage of the string method
	Error     MessageType = "error"
	Undefined MessageType = "und"
)

// This will be our DTO
type Message struct {
	Type      MessageType `json:"type"`
	Username  string      `json:"username"`
	Content   string      `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewErr(username string, err error) Message {
	return Message{Type: Error,
		Username:  username,
		Content:   err.Error(),
		Timestamp: time.Now()}
}

func (m *Message) Encode() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	return append(data, []byte("\n\n")...), nil
}

func (m *Message) Decode(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *Message) ParseContent() (MessageType, string) {
	joinParser := regexp.MustCompile(`^/join\s+(.+)$`)
	leaveParser := regexp.MustCompile(`^/leave\s*$`)

	switch m.Type {
	case Join:
		matches := joinParser.FindStringSubmatch(m.Content)
		if len(matches) == 2 {
			return m.Type, strings.TrimSpace(matches[1])
		}
		return Undefined, ""

	case Leave:
		if leaveParser.MatchString(m.Content) {
			return m.Type, ""
		}
		return Undefined, ""

	case Text, InitUser:
		return m.Type, strings.TrimSpace(m.Content)

	default:
		return Undefined, ""
	}
}
