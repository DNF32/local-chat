package main

import (
	"fmt"
	"local-chat/internal/server"
)

func main() {
	var id int = 1

	s := server.Server{}
	s.Start()
	for {
		conn, err := s.Listener.AcceptTCP()
		if err != nil {
			return
		}

		user, err := server.InitUser(conn, id)
		if err != nil {
			// Failed to init this user continue in the loop to handle a new connection
			continue
		}
		id++

		userSend := make(chan server.Event, 20)
		userRecv := make(chan server.Event, 20)
		s.SendChannels[user.ID] = userSend
		s.RecvChannels[user.ID] = userRecv
		fmt.Printf("Created send/recv channels for user %d\n", user.ID)

		go s.HandleConn(conn, user)
		fmt.Printf("Started goroutine to handle connection for user %d\n", user.ID)
	}
}
