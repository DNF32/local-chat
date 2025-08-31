package main

import (
	"errors"
	"local-chat/internal/connected_user"
	"local-chat/internal/network"
	"local-chat/internal/protocol"
	"local-chat/internal/server"
	"local-chat/internal/transport/serde"
	"log/slog"
	"net"
)

func handleClientLogin(conn *net.TCPConn, s *serde.Serde, authRepo connected_user.AuthRepository, userRepo connected_user.UserRepository, logger *slog.Logger) *connected_user.User {
	var user *connected_user.User
	var err error

	for {
		user, err = server.HandleLogin(conn, s, authRepo, userRepo)
		if err != nil {
			var intErr serde.SerdeError
			if errors.As(err, &intErr) {
				// Serialization error: retry login
				logger.Info("Serde error, retrying login:", "err", err)
				continue
			} else {
				// Other errors (e.g., invalid credentials)
				logger.Info("Login failed:", "err", err)
				errMsg := protocol.NewFailedLoginResponse(err)
				data, _ := s.EncodeEncrypted(&errMsg)
				_ = network.WriteProtocol(conn, data)
				continue
			}
		}

		// Successful login - send success response and break
		logger.Info("Login Msg", "User:", user.Username)
		successMsg := protocol.NewSucessLoginResponse(user.Username)
		data, _ := s.EncodeEncrypted(&successMsg)
		_ = network.WriteProtocol(conn, data)
		break
	}
	return user
}
