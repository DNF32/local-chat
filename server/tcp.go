package main

import (
	"fmt"
	"local-chat/models"
	"net"
	"time"
)

type Server struct {
	sendChannels map[int]chan string // for sending TO users
	recvChannels map[int]chan string // for receiving FROM users
}

func handleConn(user models.User, server Server) {
	buf := make([]byte, 1024)
	for {
		select {
		case msg := <-server.recvChannels[user.ID]:
			for Id, sendChannel := range server.sendChannels {
				if Id != user.ID {
					sendChannel <- msg
				}
			}
		case msg := <-server.sendChannels[user.ID]:
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
	sendChannels := make(map[int]chan string)
	recvChannels := make(map[int]chan string)
	server := Server{sendChannels: sendChannels,
		recvChannels: recvChannels}

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
		user := User{id: id,
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

func InitUserSessionServerSide(conn * net.TCPConn)
