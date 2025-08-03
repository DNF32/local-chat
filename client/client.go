package main

import (
	"encoding/json"
	"fmt"
	"local-chat/models"
	"local-chat/network"
	"net"
	"os"
	"time"
)

// In this modules we should do the initing of the client side strutures needed to interact with our server
// for instance we would need

//

type ChatClient struct {
	Conn     *net.TCPConn
	Incoming chan models.Message //TODO: Maybe change the type
	Outgoing chan models.Message
}

func InitID(conn *net.TCPConn, name string) (int, error) {
	m := models.Message{Type: models.InitUser,
		Username:  name,
		Content:   "None",
		Timestamp: time.Now()}

	// Send the init message
	data, err := m.Encode()
	if err != nil {
		return -1, err
	}
	_, err = conn.Write(data)
	if err != nil {
		return -1, err
	}

	// Read the response from server
	data, err = network.ReadProtocol(conn)
	if err != nil {
		return -1, err
	}

	var id int
	err = json.Unmarshal(data, &id)

	if err != nil {
		return -1, fmt.Errorf("Faild to Unmarshal the id sent by the server")
	}
	return id, nil
}

func InitUserSession() (*ChatClient, *models.User, error) {
	if len(os.Args) <= 1 {
		panic("Failed to provide Username")
	}

	name := os.Args[1]

	conn, err := net.Dial("tcp4", "localhost:8088")
	tcpConn := conn.(*net.TCPConn)
	if err != nil {
		return nil, nil, err
	}
	id, err := InitID(tcpConn, name)
	if err != nil {
		panic(fmt.Errorf("Failed to InitID: %w", err))
	}

	user := models.User{Username: name, ID: id}

	incoming := make(chan models.Message)
	outgoing := make(chan models.Message)
	client := ChatClient{Conn: tcpConn, Incoming: incoming, Outgoing: outgoing}

	// Inits the network layer between client/ui <-> server
	go network.HandleInput(tcpConn, client.Outgoing)
	go network.HandleOutput(tcpConn, client.Incoming)

	return &client, &user, nil
}
