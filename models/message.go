package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type MessageType string

const (
	InitUser MessageType = "initUser"
	Join     MessageType = "join"
	Text     MessageType = "text"
	Leave    MessageType = "leave"
)

type Message struct {
	Type      MessageType `json:"type"`
	Username  string      `json:"username"`
	Content   string      `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
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
