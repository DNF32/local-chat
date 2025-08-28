package client

import (
	"errors"
	"local-chat/internal/connected_user"
	"local-chat/internal/logger"
	"local-chat/internal/message"
	"local-chat/internal/network"
	"local-chat/internal/protocol"
	"local-chat/internal/transport"
	"local-chat/internal/transport/crypto"
	"local-chat/internal/transport/serde"
	"log/slog"
	"net"
	"time"
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

	cs := crypto.Base64Encoder{}
	sd := serde.New(&cs, logger)

	cf := ClientInfra{ // Consistent variable naming
		Logger: logger,
		conn:   tcpConn, // Use the cast connection
		serde:  sd,
	}

	return cf, nil
}

type ChatClient struct {
	Incoming chan protocol.ClientMessage
	Outgoing chan protocol.ClientMessage
	Ack      chan any

	User connected_user.User
	ClientInfra
}

func NewChatClient() (*ChatClient, error) {
	infra, err := NewClientInfra()
	if err != nil {
		return nil, err // Just return the error, don't check for specific type
	}

	// Start TCP
	incoming := make(chan protocol.ClientMessage)
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

func (c *ChatClient) Auth(auth connected_user.AuthCredentials) {
	// we will send the credential and expect a clientMessage with type initMessage if all good,

}

func sendAuthRequest(conn *net.TCPConn, auth connected_user.AuthCredentials, serde *transport.Serde) (connected_user.User, error) {

	cmsg := protocol.NewClientMessageWithAuth(auth)

	data, err := serde.EncodeEncrypted(&cmsg)
	if err != nil {
		return connected_user.User{}, err
	}

	err = network.WriteProtocol(conn, data)
	if err != nil {
		return connected_user.User{}, err
	}

	data, err = network.ReadProtocol(conn, nil)
	if err != nil {
		return connected_user.User{}, err
	}

	var u connected_user.User
	err = serde.DecodeEncrypted(data, &u)
	if err != nil {
		return connected_user.User{}, err
	}

	return u, nil
}
