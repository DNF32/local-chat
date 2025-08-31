package state

type State int

const (
	Login State = iota + 1
	Dash
	Room
)

type StateChangeMsg struct {
	NewState State
	Username string
}
