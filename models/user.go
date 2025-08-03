package models

import (
	"encoding/json"
	"fmt"
	"net"
)

type User struct {
	ID       int          `json:"id"`
	Username string       `json:"username"`
}

// This seems super familiar with a io.Reader and writer interface
// This things should have typed channels like Using the User struct and Message struct this ways this acts as a form
// of transport layer
func (u *User) Encode() ([]byte, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	return append(data, []byte("\n\n")...), nil
}

// This function parses the json data coming from ReadProtocol to the object `s`, which will be a Message or a User
func (u *User) Decode(data []byte) error {
	return json.Unmarshal(data, u)
}
