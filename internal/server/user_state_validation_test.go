package server

import (
	"local-chat/internal/message"
	"local-chat/internal/connected_user"
	"sync"
	"testing"
	"time"
)

var (
	once    sync.Once
	testSrv *Server
)

func getTestServer(t *testing.T) *Server {
	t.Helper()
	once.Do(func() {
		s := &Server{}
		s.Start()
		testSrv = s
	})
	return testSrv
}

func TestUserStateValid(t *testing.T) {
	s := getTestServer(t)

	u := connected_user.ConnectedUser{ID: 1, Username: "Tester"}
	// This user sent a valid event
	general := s.Rooms[connected_user.General]
	u.JoinRoom(general)

	if u.InRoom != true && u.Room != general {
		t.Fatal("The User state set had an issue please review .JoinRoom method")
	}

	validEvent := Event{Type: message.Text,
		Username: u.Username,
		Action:   "text",
		Content:  "this is a new message", Timestamp: time.Now()}

	newEvent, _ := s.ProcessUserEvent(&u, validEvent)
	// Check the side effect of the user we passed in
	if u.InRoom != true && u.Room != general {
		t.Fatal("ValidUserState altered the user state")
	}

	expectedEvent := Event{Type: message.Text,
		Username:  u.Username,
		Action:    "text",
		Content:   "this is a new message",
		Timestamp: validEvent.Timestamp,
		RouteRoomName:  connected_user.General}

	if equal, field := newEvent.EqualTo(&expectedEvent); equal != true {
		t.Fatalf("Expected events to be equal but got error in: %v", field)
	}
}

func TestFailingJoinEventUserInRoomValid(t *testing.T) {
	s := getTestServer(t)
	u := connected_user.ConnectedUser{ID: 1, Username: "Tester"}
	// This user sent a valid event
	general := s.Rooms[connected_user.General]
	u.JoinRoom(general)
	if !u.InRoom || u.Room != general {
		t.Fatal("The User state set had an issue please review .JoinRoom method")
	}
	validEvent := Event{Type: message.Join,
		Username: u.Username,
		Action:   "joined",
		Content:  "main", Timestamp: time.Now()}
	_, err := s.ProcessUserEvent(&u, validEvent)
	// Check the side effect of the user we passed in
	if err == nil {
		t.Fatal("Should have returned an error when user already in room")
	}

	expectedMsg := "User already in a room, to join another leave the current one: /leave"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error: '%s', got: '%s'", expectedMsg, err.Error())
	}

	if !u.InRoom || u.Room != general {
		t.Fatal("ValidUserState altered the user state")
	}
}

func TestInvalidRoomName(t *testing.T) {
	s := getTestServer(t)
	u := connected_user.ConnectedUser{ID: 1, Username: "Tester"}
	// This user sent a valid event
	// general := s.Rooms[user.General]

	validEvent := Event{Type: message.Join,
		Username: u.Username,
		Action:   "joined",
		Content:  "mained", Timestamp: time.Now()}
	_, err := s.ProcessUserEvent(&u, validEvent)
	// Check the side effect of the user we passed in
	if err == nil {
		t.Fatal("Should have returned an error for invalida RoomName")
	}

	expectedMsg := "invalid room entered: mained"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error: '%s', got: '%s'", expectedMsg, err.Error())
	}

	if u.InRoom {
		t.Fatal("ValidUserState altered the user state")
	}
}

func TestValidJoiningRoom(t *testing.T) {
	s := getTestServer(t)
	u := connected_user.ConnectedUser{ID: 1, Username: "Tester"}
	general := s.Rooms[connected_user.General]

	validEvent := Event{Type: message.Join,
		Username: u.Username,
		Action:   "joined",
		Content:  "general", Timestamp: time.Now()}
	newEvent, err := s.ProcessUserEvent(&u, validEvent)
	// Check the side effect of the user we passed in
	if err != nil {
		t.Fatal("Got and unexpected error during ValidUserState")
	}

	expectedEvent := Event{Type: message.Join,
		Username: u.Username,
		Action:   "joined",
		Content:  "general", Timestamp: time.Now(),
		RouteRoomName: connected_user.General,
	}

	// We need to the new user state and the expected event if matched we only need to see RoomName
	if expectedEvent.RouteRoomName != newEvent.RouteRoomName {
		t.Fatalf("Expected error: '%s', got: '%s'", expectedEvent.RouteRoomName, expectedEvent.RouteRoomName)
	}
	if !u.InRoom || u.Room != general {
		t.Errorf("Expected error: '%v', got: '%v'", u.InRoom, true)
		t.Fatalf("Expected error: '%s', got: '%s'", u.Room.RoomName, "general")
	}
}

func TestFailingTextEventUserNotInRoom(t *testing.T) {
	s := getTestServer(t)
	u := connected_user.ConnectedUser{ID: 1, Username: "Tester"}

	// Verify user is not in any room initially
	if u.InRoom {
		t.Fatal("User should not be in a room initially")
	}

	validEvent := Event{Type: message.Text,
		Username:  u.Username,
		Action:    "",
		Content:   "Hello everyone!",
		Timestamp: time.Now()}

	_, err := s.ProcessUserEvent(&u, validEvent)

	// Should return an error
	if err == nil {
		t.Fatal("Should have returned an error when user not in room tries to send text")
	}

	expectedMsg := "User not in room, so can't send a message"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error: '%s', got: '%s'", expectedMsg, err.Error())
	}

	// Verify user state wasn't altered
	if u.InRoom {
		t.Fatal("ValidUserState should not have altered user state")
	}
}

func TestTextEventUserInRoom(t *testing.T) {
	s := getTestServer(t)
	u := connected_user.ConnectedUser{ID: 1, Username: "Tester"}
	general := s.Rooms[connected_user.General]
	u.JoinRoom(general)

	// Verify user IS in a room after joining
	if !u.InRoom {
		t.Fatal("User should be in a room after joining")
	}

	validEvent := Event{Type: message.Text,
		Username:  u.Username,
		Action:    "",
		Content:   "Hello everyone!",
		Timestamp: time.Now()}

	newEvent, err := s.ProcessUserEvent(&u, validEvent)

	// Should NOT return an error
	if err != nil {
		t.Fatalf("Should not have returned an error, got: %v", err)
	}

	expectedEvent := Event{Type: message.Text,
		Username:  u.Username,
		Action:    "",
		Content:   "Hello everyone!",
		Timestamp: time.Now(),
		RouteRoomName:  connected_user.General}

	if expectedEvent.RouteRoomName != newEvent.RouteRoomName {
		t.Fatalf("Expected RoomName: '%s', got: '%s'", expectedEvent.RouteRoomName, newEvent.RouteRoomName)
	}

	if !u.InRoom || u.Room != general {
		t.Errorf("Expected user to be in room: %v", u.InRoom)
		t.Fatalf("Expected user room: '%s', got: '%s'", general.RoomName, u.Room.RoomName)
	}
}

func TestLeaveEventUserInRoom(t *testing.T) {
	s := getTestServer(t)
	u := connected_user.ConnectedUser{ID: 1, Username: "Tester"}
	general := s.Rooms[connected_user.General]
	u.JoinRoom(general)

	// Verify user is in room initially
	if !u.InRoom {
		t.Fatal("User should be in a room before leaving")
	}

	validEvent := Event{Type: message.Leave,
		Username:  u.Username,
		Action:    "left",
		Content:   "",
		Timestamp: time.Now()}

	newEvent, err := s.ProcessUserEvent(&u, validEvent)

	if err != nil {
		t.Fatalf("Should not have returned an error, got: %v", err)
	}

	// Should return the room name they left
	if newEvent.RouteRoomName != connected_user.General {
		t.Fatalf("Expected RoomName: '%s', got: '%s'", connected_user.General, newEvent.RouteRoomName)
	}

	// User should no longer be in room after leaving
	if u.InRoom {
		t.Fatal("User should not be in room after leaving")
	}
}

func TestLeaveEventUserNotInRoom(t *testing.T) {
	s := getTestServer(t)
	u := connected_user.ConnectedUser{ID: 1, Username: "Tester"}

	// Verify user is not in room initially
	if u.InRoom {
		t.Fatal("User should not be in a room initially")
	}

	validEvent := Event{Type: message.Leave,
		Username:  u.Username,
		Action:    "left",
		Content:   "",
		Timestamp: time.Now()}

	_, err := s.ProcessUserEvent(&u, validEvent)

	if err == nil {
		t.Fatal("Should have returned an error when user not in room tries to leave")
	}

	// Check for specific error from LeaveRoom method
	// (You'll need to adjust this based on what LeaveRoom returns)
}

func TestUnsupportedMessageType(t *testing.T) {
	s := getTestServer(t)
	u := connected_user.ConnectedUser{ID: 1, Username: "Tester"}

	// Create event with unsupported message type
	validEvent := Event{Type: message.Error, // Or whatever represents unsupported type
		Username:  u.Username,
		Action:    "",
		Content:   "test",
		Timestamp: time.Now()}

	_, err := s.ProcessUserEvent(&u, validEvent)

	if err == nil {
		t.Fatal("Should have returned an error for unsupported message type")
	}

	expectedMsg := "The other message types are not supported"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error: '%s', got: '%s'", expectedMsg, err.Error())
	}
}

// Test Event parsing

func TestValidateNewEvent(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		input       message.Message
		wantErr     bool
		wantAction  string
		wantContent string
	}{
		{
			name: "InitUser event",
			input: message.Message{
				Type:      message.InitUser,
				Username:  "system",
				Content:   "name",
				Timestamp: now,
			},
			wantErr:     false,
			wantAction:  "initUser", // adjust based on your mapping rules
			wantContent: "name",
		},
		{
			name: "Join event",
			input: message.Message{
				Type:      message.Join,
				Username:  "alice",
				Content:   "/join main",
				Timestamp: now,
			},
			wantErr:     false,
			wantAction:  "joined", // e.g. parsed action
			wantContent: "main",   // e.g. extracted room
		},
		{
			name: "Text event",
			input: message.Message{
				Type:      message.Text,
				Username:  "bob",
				Content:   "this is the end",
				Timestamp: now,
			},
			wantErr:     false,
			wantAction:  "send",
			wantContent: "this is the end",
		},
		{
			name: "Leave event",
			input: message.Message{
				Type:      message.Leave,
				Username:  "charlie",
				Content:   "/leave",
				Timestamp: now,
			},
			wantErr:     false,
			wantAction:  "left",
			wantContent: "",
		},
		{
			name: "Error event",
			input: message.Message{
				Type:      message.Error,
				Username:  "system",
				Content:   "Failed to deliver message",
				Timestamp: now,
			},
			wantErr: true,
		},
		{
			name: "Undefined event",
			input: message.Message{
				Type:      message.Undefined,
				Username:  "ghost",
				Content:   "???",
				Timestamp: now,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateNewEvent(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateNewEvent(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}

			if !tt.wantErr {
				if got.Action != tt.wantAction {
					t.Errorf("Action = %v, want %v", got.Action, tt.wantAction)
				}
				if got.Content != tt.wantContent {
					t.Errorf("Content = %v, want %v", got.Content, tt.wantContent)
				}
			}
		})
	}
}
