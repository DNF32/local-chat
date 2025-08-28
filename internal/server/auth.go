package server

import (
	"errors"
	"fmt"
	"local-chat/internal/connected_user"
	"local-chat/internal/network"
	"local-chat/internal/protocol"
	"local-chat/internal/transport/serde"
	"net"
)

var ErrInvalidMsgType = errors.New("Remote sent ")

func HandleLogin(
	conn *net.TCPConn,
	serde *serde.Serde,
	authRepo connected_user.AuthRepository,
	userRepo connected_user.UserRepository,
) (*connected_user.User, error) {

	data, err := network.ReadProtocol(conn, nil)
	if err != nil {
		return nil, err
	}

	var cmsg protocol.ClientMessage
	err = serde.DecodeEncrypted(data, &cmsg)
	if err != nil {
		return nil, err
	}

	if cmsg.Type != protocol.EventTypeLogin {
		return nil, fmt.Errorf(
			"%v sent wrong message type %s, when only %v is admitted",
			conn.RemoteAddr().String(),
			cmsg.Type,
			protocol.EventTypeLogin,
		)
	}

	user, err := connected_user.Login(cmsg.AuthCredentials, authRepo, userRepo)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
