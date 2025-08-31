package protocol

import (
	"encoding/json"
	"local-chat/internal/connected_user"
)

type ClientMessage struct {
	Type     EventType `json:"type"` // join, leave, chat, login, etc.
	Username string    `json:"username"`
	Content  string    `json:"content"` // room name for join/leave, message for chat
	connected_user.AuthCredentials
}

func NewClientMessageWithAuth(auth connected_user.AuthCredentials) ClientMessage {
	return ClientMessage{
		Type:            EventTypeLogin,
		AuthCredentials: auth,
	}
}

func NewClientJoinMessage(username, roomName string) ClientMessage {
	return ClientMessage{
		Type:     EventTypeJoin,
		Username: username,
		Content:  roomName,
	}
}

func NewClientLeaveMessage(username string) ClientMessage {
	return ClientMessage{
		Type:     EventTypeLeave,
		Username: username,
	}
}

func NewClientTextMessage(username, content string) ClientMessage {
	return ClientMessage{
		Type:     EventTypeLeave,
		Username: username,
		Content:  content,
	}
}

func (cm *ClientMessage) Encode() ([]byte, error) {
	return json.Marshal(cm)
}

func (cm *ClientMessage) Decode(data []byte) error {
	return json.Unmarshal(data, cm)
}
