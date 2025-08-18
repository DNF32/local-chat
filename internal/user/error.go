package user

import "fmt"

type InvalidRoomError struct {
	Room string
}

func (e *InvalidRoomError) Error() string {
	return fmt.Sprintf("invalid room entered: %s", e.Room)
}
