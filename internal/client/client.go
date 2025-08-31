package client

import (
	"errors"
	"fmt"
	"local-chat/internal/connected_user"
	"local-chat/internal/logger"
	"local-chat/internal/protocol"
	"local-chat/internal/transport/cryptol"
	"local-chat/internal/transport/serde"
	"log/slog"
	"net"
)

var ErrInvalidCredential = errors.New("Invalid credentials, please try again.")
var ErrReachingServer = errors.New("Failed to reach server at localhost:8088")

type ClientInfra struct {
	conn   *net.TCPConn
	serde  *serde.Serde
	Logger *slog.Logger
}

func NewClientInfra() (ClientInfra, error) {
	conn, err := net.Dial("tcp4", "localhost:8088")
	if err != nil {
		return ClientInfra{}, ErrReachingServer // Check error BEFORE casting
	}

	tcpConn := conn.(*net.TCPConn) // Cast after confirming no error

	logger, err := logger.NewFileLogger(logger.CLIENT_INFRA_LOG_PATH)
	if err != nil {
		conn.Close() // Clean up connection on logger error
		return ClientInfra{}, err
	}

	cs := cryptol.Base64Encoder{}
	sd := serde.New(&cs, logger)

	cf := ClientInfra{ // Consistent variable naming
		Logger: logger,
		conn:   tcpConn, // Use the cast connection
		serde:  sd,
	}

	return cf, nil
}

type ChatClient struct {
	Incoming chan protocol.ServerResponse
	Outgoing chan protocol.ClientMessage
	Ack      chan any

	User connected_user.User
	ClientInfra
}

func (cc ChatClient) HandleInput() {
	HandleInput(cc.ClientInfra.conn, cc.Outgoing, cc.ClientInfra.serde, cc.ClientInfra.Logger)

}

func (cc ChatClient) HandleOutput() {
	HandleOutput(cc.ClientInfra.conn, cc.Incoming, cc.ClientInfra.serde, cc.ClientInfra.Logger)
}

func NewChatClient() (*ChatClient, error) {
	infra, err := NewClientInfra()
	if err != nil {
		return nil, err // Just return the error, don't check for specific type
	}

	// Start TCP
	incoming := make(chan protocol.ServerResponse)
	outgoing := make(chan protocol.ClientMessage)
	ack := make(chan any)

	c := ChatClient{
		Incoming:    incoming,
		Outgoing:    outgoing,
		Ack:         ack,
		ClientInfra: infra,
	}

	return &c, nil // Return pointer to ChatClient
}

type FailedLoginErr struct {
	serverError error
}

func (err FailedLoginErr) Error() string {
	return fmt.Sprintf("Failed login: %v", err.serverError.Error())
}

type Username string

func (cc *ChatClient) SendAuthRequest(username, password string) {
	auth := connected_user.NewAuthCredentials(username, password)
	cmsg := protocol.NewClientMessageWithAuth(auth)

	cc.Outgoing <- cmsg
}
