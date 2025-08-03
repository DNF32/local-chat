package models

import "time"

type MessageType string

const (
	Join  MessageType = "join"
	Text  MessageType = "text"
	Leave MessageType = "leave"
)

type Message struct {
	Type      string    `json:"type"`
	Username  string    `json:"username"`
	content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
