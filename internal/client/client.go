package client

import (
	"fmt"
	"local-chat/internal/message"
	"local-chat/internal/network"
	"local-chat/internal/user"
	"log/slog"
	"net"
	"os"
	"time"
)

// In this modules we should do the initing of the client side strutures needed to interact with our server
// for instance we would need

//

type ChatClient struct {
	Conn     *net.TCPConn
	Incoming chan message.Message //TODO: Maybe change the type
	Outgoing chan message.Message
}

func InitUser(conn *net.TCPConn, name string) (user.User, error) {
	fmt.Printf("InitID: started for user %q\n", name)

	m := message.Message{
		Type:      message.InitUser,
		Username:  name,
		Content:   "None",
		Timestamp: time.Now(),
	}

	data, err := m.Encode()
	if err != nil {
		fmt.Printf("InitID: failed to encode message: %v\n", err)
		return user.User{}, err
	}
	fmt.Printf("InitID: message encoded, sending %d bytes\n", len(data))

	fmt.Printf("InitID: message encoded,%s\n", data)
	n, err := conn.Write(data)
	if err != nil {
		fmt.Printf("InitID: failed to write to connection: %v\n", err)
		return user.User{}, err
	}
	fmt.Printf("InitID: wrote %d bytes to connection\n", n)

	data, err = network.ReadProtocol(conn, nil)
	if err != nil {
		fmt.Printf("InitID: failed to read response: %v\n", err)
		return user.User{}, err
	}
	fmt.Printf("InitID: received %d bytes from server\n", len(data))

	var u user.User
	err = u.Decode(data)
	if err != nil {
		fmt.Printf("InitID: failed to unmarshal response: %v\n", err)
		return user.User{}, fmt.Errorf("failed to unmarshal ID: %w", err)
	}

	return u, nil
}

func InitUserSession(logger *slog.Logger) (*ChatClient, *user.User, error) {
	if len(os.Args) <= 1 {
		panic("Failed to provide Username")
	}

	name := os.Args[1]

	conn, err := net.Dial("tcp4", "localhost:8088")
	tcpConn := conn.(*net.TCPConn)
	if err != nil {
		return nil, nil, err
	}

	user, err := InitUser(tcpConn, name)
	if err != nil {
		panic(fmt.Errorf("Failed to InitID: %w", err))
	}
	fmt.Println("ID is", user.ID)

	incoming := make(chan message.Message)
	outgoing := make(chan message.Message)
	client := ChatClient{Conn: tcpConn, Incoming: incoming, Outgoing: outgoing}

	// Inits the network layer between client/ui <-> server
	go network.HandleInput(tcpConn, client.Outgoing, logger)
	go network.HandleOutput(tcpConn, client.Incoming, logger)

	return &client, &user, nil
}
