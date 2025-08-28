package protocol

import (
	"encoding/json"
	"local-chat/internal/connected_user"
	"time"
)

type EventType string

const (
	// Client-initiated actions
	EventTypeJoin   EventType = "join"
	EventTypeLeave  EventType = "leave"
	EventTypeChat   EventType = "chat"
	EventTypeLogin  EventType = "login"
	EventTypeLogout EventType = "logout"

	// Server responses
	EventTypeAck         EventType = "ack"
	EventTypeError       EventType = "error"
	EventTypeFailedLogin EventType = "FailedLogin"
	EventTypeSucessLogin EventType = "FailedLogin"

	// Server broadcasts
	EventTypeUserJoined    EventType = "user_joined"
	EventTypeUserLeft      EventType = "user_left"
	EventTypeChatBroadcast EventType = "chat_broadcast"
)

type Event struct {
	ClientMessage

	TypeMetaData  map[string]any          // Data associated with the content sent
	RouteRoomName connected_user.RoomName // Room to route this message
	Timestamp     time.Time               `json:"timestamp"`
}

func EventFromMessage(msg ClientMessage, metadata map[string]any, room connected_user.RoomName) Event {
	return Event{
		ClientMessage: msg,
		TypeMetaData:  metadata,
		RouteRoomName: room,
		Timestamp:     time.Now(),
	}
}

/*
	This function assumed that the content comes in a form that can be consumed by the backend

For instance user sends: /join main

Event.Content = "main" and not "/join main"

Note: assume input is cleaned of whitespace
*/
func (e *Event) ContentValidation() error {
	switch e.Type {
	case EventTypeJoin:
		roomName := e.Content
		roomNameValidated, err := connected_user.ValidateRoom(roomName)
		if err != nil {
			return err
		}
		e.Content = string(roomNameValidated)
	default:
	}
	return nil
}

func (e *Event) Encode() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) Decode(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e Event) ToResponse() ServerResponse {
	sr := ServerResponse{}
	// we need to fill in the basic stuff
	sr.Username = e.Username
	sr.RouteRoomName = e.RouteRoomName

	switch e.Type {
	case EventTypeJoin:
		e.Type = EventTypeUserJoined
	case EventTypeLeave:
		e.Type = EventTypeUserLeft
	case EventTypeChat:
		e.Type = EventTypeChatBroadcast
		sr.Content = e.Content // if it's chat this is the message, otherwise this isn't needed
	}
	return sr
}
