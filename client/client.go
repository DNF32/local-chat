package main

import (
	"local-chat/network"
	"net"
)

// In this modules we should do the initing of the client side strutures needed to interact with our server
// for instance we would need

//

type ChatClient struct {
	Conn     *net.TCPConn
	Incoming chan string //TODO: Maybe change the type
	Outgoing chan string
}

func InitChatClient() (*ChatClient, error) {
	conn, err := net.Dial("tcp4", "localhost:8088")
	tcpConn := conn.(*net.TCPConn)
	if err != nil {
		return nil, err
	}

	incoming := make(chan string)
	outgoing := make(chan string)
	client := ChatClient{Conn: tcpConn, Incoming: incoming, Outgoing: outgoing}

	go network.HandleInput(tcpConn, client.Outgoing)
	go network.HandleOutput(tcpConn, client.Incoming)

	return &client, nil
}
