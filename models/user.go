package models

import "net"

type User struct {
	ID       int          `json:"id"`
	Username string       `json:"username"`
	Conn     *net.TCPConn `json:"-"`
}
