package main

import (
	"encoding/json"
	"fmt"
	"local-chat/models"
	"local-chat/network"
	"net"
	"time"
)

type Server struct {
	Users        map[int]models.User
	SendChannels map[int]chan models.Message // for sending TO users
	RecvChannels map[int]chan models.Message // for receiving FROM users
}

// This enables the command list active users
// we could to something even better send a ping through the connection to see if the user is still live

func handleConn(conn *net.TCPConn, user models.User, server Server) {
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
			data, err := msg.Encode()
			if err != nil {
				return
			}
			for len(data) > 0 {
				n, err := conn.Write(data)
				if err != nil {
					return
				}
				data = data[n:]
			}
		default:
			msgBytes, err := ParseMsgSent(conn, buf)
			if err != nil {
				return
			}

			var msg models.Message
			err = msg.Decode(msgBytes)
			if err != nil {
				return
			}

			switch msg.Type {
			case models.Join:
				// Add the user to the room general
				server.Users[user.ID] = user
				server.RecvChannels[user.ID] <- msg
			case models.Leave:
				// Remove the user from the active list
				delete(server.Users, user.ID)
				server.RecvChannels[user.ID] <- msg
			case models.Text:
				// Broadcast the message
				server.RecvChannels[user.ID] <- msg
			default:
				panic("Incorrect message type received on server")
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

		user := models.User{ID: id,
			Username: initMsg.Username,
		}
		id += 1
		data, _ = json.Marshal(user.ID)

		for len(data) > 0 {
			n, err := conn.Write(data)
			if err != nil {
				return
			}
			data = data[n:]
		}

		userSend := make(chan models.Message, 20)
		userRecv := make(chan models.Message, 20)
		server.SendChannels[user.ID] = userSend
		server.RecvChannels[user.ID] = userRecv

		go handleConn(conn, user, server)
	}

}

func InitUserSessionServerSide(conn *net.TCPConn)
