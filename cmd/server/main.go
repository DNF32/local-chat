package main

import (
	"fmt"
	"local-chat/internal/connected_user"
	"local-chat/internal/protocol"
	"local-chat/internal/server"
	"net"
)

func main() {

	s := server.Server{}
	authRepo, _ := connected_user.NewSQLiteAuthRepo("")
	userRepo, _ := connected_user.NewSQLiteUserRepo("")
	s.Start()

	for {
		conn, err := s.Listener.AcceptTCP()
		if err != nil {
			continue // don't exit the server on accept error
		}

		go func(conn *net.TCPConn) {
			defer conn.Close() // make sure connection is closed when done

			var user *connected_user.User
			user = handleClientLogin(conn, s.Serde, authRepo, userRepo, s.Logger)
			connected_user := connected_user.ConnectedUser{User: *user, InRoom: false, Room: nil}

			userSend := make(chan protocol.Event, 20)
			userRecv := make(chan protocol.Event, 20)
			s.SendChannels[user.ID] = userSend
			s.RecvChannels[user.ID] = userRecv
			fmt.Printf("Created send/recv channels for user %d\n", user.ID)

			s.HandleConn(conn, &connected_user) // handle the connection for this user
			fmt.Printf("Started goroutine to handle connection for user %d\n", user.ID)
		}(conn)
	}
}
