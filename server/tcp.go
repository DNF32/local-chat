package main

import (
	"fmt"
	"local-chat/models"
	"local-chat/network"
	"net"
	"time"
)

type Server struct {
	Users        []models.User
	SendChannels map[int]chan models.Message // for sending TO users
	RecvChannels map[int]chan models.Message // for receiving FROM users
}

// Need to store the user in the slice of the server, maybe we could also use a map, this is just to keep in sync, knowin the users would be useful for multiple rooms
// This enables the command list active users
// we could to something even better send a ping through the connection to see if the user is still live

func handleConn(user models.User, server Server) {
	buf := make([]byte, 1024)
	for {
		select {
		case msg := <-server.RecvChannels[user.ID]:
			for Id, sendChannel := range server.SendChannels {
				if Id != user.ID {
					sendChannel <- msg
				}
			}
		case msg := <-server.SendChannels[user.ID]:
			msg = msg + "\n\n"
			data := []byte(msg)
			for len(data) > 0 {
				n, err := user.Conn.Write(data)
				if err != nil {
					return
				}
				data = data[n:]
			}
		default:
			msgBytes, err := ParseMsgSent(user.Conn, buf)
			if err != nil {
				return
			}
			if len(msgBytes) > 0 {
				server.recvChannels[user.ID] <- string(msgBytes)
			}
		}
	}
}

func ParseMsgSent(conn *net.TCPConn, buf []byte) ([]byte, error) {
	var msgBytes []byte
	for {
		conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
		n, err := conn.Read(buf)

		if n > 0 {
			msgBytes = append(msgBytes, buf[:n]...)
			if len(msgBytes) >= 2 &&
				msgBytes[len(msgBytes)-2] == '\n' &&
				msgBytes[len(msgBytes)-1] == '\n' {
				break
			}
		}

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Timeout is expected, break out of inner loop and try select again
				break
			}
			// Real error
			return nil, err
		}
	}
	return msgBytes, nil
}
func main() {
	var id int = 1
	intBuf := make([]byte, 0, 64)
	sendChannels := make(map[int]chan models.Message)
	recvChannels := make(map[int]chan models.Message)
	server := Server{SendChannels: sendChannels,
		RecvChannels: recvChannels}

	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8088")
	if err != nil {
		fmt.Println("Error resolving address")
		panic(err)
	}
	listen, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		fmt.Println("Error resolving address")
		panic(err)
	}
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("Error resolving address")
			panic(err)
		}

		data, err := network.ReadProtocol(conn, intBuf)
		if err != nil {
			fmt.Println("Error reading init message")
			panic(err)
		}

		var initMsg models.Message
		err = initMsg.Decode(data)
		if err != nil {
			panic(err)
		}

		user := User{ID: id,
			name: "Adele",
			conn: conn}
		id += 1

		userSend := make(chan string, 20)
		userRecv := make(chan string, 20)
		server.sendChannels[user.id] = userSend
		server.recvChannels[user.id] = userRecv

		go handleConn(user, server)
	}

}

func InitUserSessionServerSide(conn *net.TCPConn)
