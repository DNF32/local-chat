package connected_user

import (
	"encoding/json"
	"fmt"
	"strings"
)

type RoomName string

const (
	General RoomName = "general"
	Main    RoomName = "main"
	Fitness RoomName = "fitness"
)

type Room struct {
	RoomName    RoomName
	ActiveUsers map[int]struct{} // maybe just use an array of ints because this will be the user ids
}

func NewRoom(name RoomName) *Room {
	return &Room{RoomName: name,
		ActiveUsers: make(map[int]struct{})}
}

func ValidateRoom(room string) (RoomName, error) {
	// Just doing some normalization
	nroom := strings.TrimSpace(room)
	nroom = strings.ToLower(nroom)
	switch RoomName(nroom) {
	case General, Main, Fitness:
		return RoomName(nroom), nil
	default:
		return RoomName(""), fmt.Errorf("invalid room entered: %s", room)
	}
}

type User struct {
	ID       int64    `json:"id"`
	Username string `json:"username"`
}

type ConnectedUser struct {
	User
	InRoom bool `json:"inroom"`
	Room   *Room
}

func (u *ConnectedUser) LeaveRoom() (RoomName, error) {
	if !u.InRoom || u.Room == nil {
		return RoomName(""), fmt.Errorf("Can't leave room, current user isn't in one")
	}
	roomName := u.Room.RoomName
	delete(u.Room.ActiveUsers, u.ID)
	u.InRoom = false
	u.Room = nil
	return roomName, nil
}
func (u *ConnectedUser) JoinRoom(room *Room) error {
	if u.InRoom {
		return fmt.Errorf("User already in a room, to join another leave the current one: /leave")
	}
	u.InRoom = true
	u.Room = room
	room.ActiveUsers[u.ID] = struct{}{}
	return nil
}

// This seems super familiar with a io.Reader and writer interface
// This things should have typed channels like Using the User struct and Message struct this ways this acts as a form
// of transport layer
func (u *ConnectedUser) Encode() ([]byte, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	return data, nil
}
func (u *ConnectedUser) Decode(data []byte) error {
	return json.Unmarshal(data, u)
}
