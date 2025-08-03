package models

import (
	"local-chat/network"
	"net"
)

type User struct {
	ID       int          `json:"id"`
	Username string       `json:"username"`
	Conn     *net.TCPConn `json:"-"`
}

func (m *User) Enconde() ([]byte, error) {
	data, err := network.Enconde(m)
	if err != nil {
		return nil, err
	}
	return data, nil
}
