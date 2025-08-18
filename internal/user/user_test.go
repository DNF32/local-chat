package user

import (
	"fmt"
	"testing"
)

func TestNewRoomValid(t *testing.T) {
	tests := []struct {
		Input        RoomName
		ExpectedName string
	}{
		{Input: General, ExpectedName: "general"},
		{Input: Main, ExpectedName: "main"},
		{Input: Fitness, ExpectedName: "fitness"},
	}

	for _, test := range tests {
		r := NewRoom(test.Input)
		if r.RoomName != test.ExpectedName {
			t.Errorf("got %v, want %v", r.RoomName, test.ExpectedName)
		}

	}
}

func TestValidRoomName(t *testing.T) {
	tests := []struct {
		Input        string
		ExpectedName error
	}{
		{Input: string(General), ExpectedName: nil},
		{Input: string(Main), ExpectedName: nil},
		{Input: string(Fitness), ExpectedName: nil},
	}

	for _, test := range tests {
		err := ValidateRoom(test.Input)
		if err != test.ExpectedName {
			t.Errorf("got %v, want %v", err.Error(), test.ExpectedName.Error())
		}

	}
}

func TestInvalidRoomName(t *testing.T) {
	tests := []struct {
		Input        string
		ExpectedName error
	}{
		{Input: "this", ExpectedName: fmt.Errorf("invalid room: %s", "this")},
		{Input: "Not a valida name", ExpectedName: fmt.Errorf("invalid room: %s", "Not a valida name")},
	}

	for _, test := range tests {
		err := ValidateRoom(test.Input)
		if err.Error() != test.ExpectedName.Error() {
			t.Errorf("got %v, want %v", err.Error(), test.ExpectedName.Error())
		}
	}
}

func TestLeavingInitUser(t *testing.T) {
	u := User{}

	err := u.LeaveRoom()
	if err != nil {
		//fmt.Printf("the init value of InRoom is: %v\n", u.InRoom) so this
		fmt.Println(err.Error())
	}
}

func TestLeavingUserValid(t *testing.T) {
	room := NewRoom(General)
	u := User{}

	err := u.JoinRoom(room)
	if err != nil {
		t.Fatal("Got: ", err)
	}
	err = u.LeaveRoom()
	if err != nil {
		t.Fatal("Got: ", err)
	}
}

func TestInvalidLeave(t *testing.T) {
	room := NewRoom(General)
	u := User{}

	err := u.JoinRoom(room)
	_ = u.LeaveRoom()
	err = u.LeaveRoom()
	if err == nil {
		t.Fatal("Got a nil err when leaving a a user not in a room ")
	}
}

func TestJoiningRoom(t *testing.T) {

	room := NewRoom(General)
	u := User{ID: 1}

	_ = u.JoinRoom(room)

	user := room.ActiveUsers[1]
	if user != &u {
		t.Fatal("User not active in room but we joined it")
	}

}

func TestTryingToJoinTwoRooms(t *testing.T) {

	room1 := NewRoom(General)
	room2 := NewRoom(Main)
	u := User{ID: 1}

	_ = u.JoinRoom(room1)

	err := u.JoinRoom(room2)
	if err == nil {
		t.Error("Expected error when joining room while already in another room")
	}

}

func TestJoinAndLeaveRoom(t *testing.T) {
	room := NewRoom(General)
	u := User{ID: 1, Username: "testuser"}

	// Test joining room
	err := u.JoinRoom(room)
	if err != nil {
		t.Fatal("Failed to join room:", err)
	}

	// Verify user is in room
	if !u.InRoom {
		t.Error("User should be in room after joining")
	}
	if u.Room != room {
		t.Error("User's room should be set to the joined room")
	}
	if room.ActiveUsers[u.ID] != &u {
		t.Error("User should be in room's ActiveUsers map")
	}

	// Test leaving room
	err = u.LeaveRoom()
	if err != nil {
		t.Fatal("Failed to leave room:", err)
	}

	// Verify user left room
	if u.InRoom {
		t.Error("User should not be in room after leaving")
	}
	if u.Room != nil {
		t.Error("User's room should be nil after leaving")
	}
	if _, exists := room.ActiveUsers[u.ID]; exists {
		t.Error("User should not be in room's ActiveUsers map after leaving")
	}
}
