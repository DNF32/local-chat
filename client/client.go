package main

import (
	"local-chat/models"
	"local-chat/network"
	"net"
	"os"
)

// In this modules we should do the initing of the client side strutures needed to interact with our server
// for instance we would need

//

type ChatClient struct {
	Incoming chan string //TODO: Maybe change the type
	Outgoing chan string
}

func InitUserSession() (*ChatClient, *models.User, error) {
	if len(os.Args) <= 1 {
		panic("Failed to provide Username")
	}

	name = os.Args[1]

	conn, err := net.Dial("tcp4", "localhost:8088")
	tcpConn := conn.(*net.TCPConn)
	if err != nil {
		return nil, err
	}

	user = models.User{Conn: conn, ID: 1}

	incoming := make(chan string)
	outgoing := make(chan string)
	client := ChatClient{Conn: tcpConn, Incoming: incoming, Outgoing: outgoing}

	go network.HandleInput(tcpConn, client.Outgoing)
	go network.HandleOutput(tcpConn, client.Incoming)

	return &client, &user nil
}
