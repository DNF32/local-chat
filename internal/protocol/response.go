package protocol

import (
	"encoding/json"
	"local-chat/internal/connected_user"
	"time"
)

type ServerResponse struct {
	Type          EventType               `json:"type"`        // ack, error, chat_broadcast, etc.
	ActionType    EventType               `json:"action_type"` // This is used in case of ack is sent, we need to say to the user if it's all ok or not depending
	Username      string                  `json:"username,omitempty"`
	Content       string                  `json:"content,omitempty"`
	RouteRoomName connected_user.RoomName `json:"room_name,omitempty"`
	//Success   bool                    `json:"success,omitempty"`
	//ErrorCode string    `json:"error_code,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func NewACKResponse(event Event) ServerResponse {
	return ServerResponse{
		Type:          EventTypeAck,
		Username:      event.Username,
		RouteRoomName: event.RouteRoomName,
		ActionType:    event.Type,
		Content:       event.Content,
		Timestamp:     time.Now(), // you might want to set this
	}
}

func NewErrResponse(event Event, err error) ServerResponse {
	return ServerResponse{
		Type:          EventTypeError,
		Username:      event.Username,
		RouteRoomName: event.RouteRoomName,
		ActionType:    event.Type,
		Content:       err.Error(),
		Timestamp:     time.Now(), // you might want to set this
	}
}

func NewFailedLoginResponse(err error) ServerResponse {
	return ServerResponse{
		Type:      EventTypeFailedLogin,
		Content:   err.Error(),
		Timestamp: time.Now(),
	}

}

func NewSucessLoginResponse(username string) ServerResponse {
	return ServerResponse{
		Type:      EventTypeSucessLogin,
		Username:  username,
		Timestamp: time.Now(),
	}

}

func (sr *ServerResponse) Encode() ([]byte, error) {
	return json.Marshal(sr)
}

func (sr *ServerResponse) Decode(data []byte) error {
	return json.Unmarshal(data, sr)
}
