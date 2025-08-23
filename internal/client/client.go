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

func InitUser(conn *net.TCPConn, name string) error {
	fmt.Printf("InitID: started for user %q\n", name)

	m := message.Message{
		Type:      message.InitUser,
		Username:  name,
		Content:   "None",
		Timestamp: time.Now(),
	}

	data, err := m.Encode()
	if err != nil {
		err = fmt.Errorf("InitID: failed to encode message: %v\n", err)
		return err
	}

	n, err := conn.Write(data)
	if err != nil {
		err = fmt.Errorf("InitID: failed to write to connection: %v\n", err)
		return err
	}

	data, err = network.ReadProtocol(conn, nil)
	if err != nil {
		err = fmt.Errorf("InitID: failed to read response: %v\n", err)
		return err
	}

	var u user.User
	err = u.Decode(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal ID: %w", err)
	}
	return nil
}

func InitUserSession(logger *slog.Logger) (*ChatClient, error) {
	if len(os.Args) <= 1 {
		panic("Failed to provide Username")
	}

	name := os.Args[1]

	conn, err := net.Dial("tcp4", "localhost:8088")
	tcpConn := conn.(*net.TCPConn)
	if err != nil {
		return nil, err
	}

	err = InitUser(tcpConn, name)
	if err != nil {
		panic(fmt.Errorf("Failed to InitID: %w", err))
	}

	incoming := make(chan message.Message)
	outgoing := make(chan message.Message)
	client := ChatClient{Conn: tcpConn, Incoming: incoming, Outgoing: outgoing}

	// Inits the network layer between client/ui <-> server
	go network.HandleInput(tcpConn, client.Outgoing, logger)
	go network.HandleOutput(tcpConn, client.Incoming, logger)

	return &client, nil
}
