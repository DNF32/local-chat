package server

import (
	"fmt"
	"log/slog"
	"net"

	"local-chat/internal/connected_user"
	"local-chat/internal/logger"
	"local-chat/internal/network"
	"local-chat/internal/protocol"
	"local-chat/internal/transport/cryptol"
	"local-chat/internal/transport/serde"
)

type Server struct {
	Users        map[int64]*connected_user.ConnectedUser
	Rooms        map[connected_user.RoomName]*connected_user.Room
	SendChannels map[int64]chan protocol.Event // for sending TO users
	RecvChannels map[int64]chan protocol.Event // for receiving FROM users

	Listener *net.TCPListener
	Serde    *serde.Serde
	Logger   *slog.Logger
}

func (s *Server) Start() {
	logger, _ := logger.NewFileLogger(logger.SERVER_LOG_PATH)
	s.Logger = logger

	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8088")
	if err != nil {
		s.Logger.Info("Error resolving address")
		panic(err)
	}
	fmt.Printf("Listening on %s\n", addr.String())

	listen, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		s.Logger.Info("Error starting TCP listener")
		panic(err)
	}

	s.Logger.Info("TCP listener started, waiting for connections...")
	s.Listener = listen

	cs := cryptol.Base64Encoder{}
	sd := serde.New(&cs, logger)
	s.Serde = sd

	sendChannels := make(map[int64]chan protocol.Event)
	recvChannels := make(map[int64]chan protocol.Event)

	s.RecvChannels = recvChannels
	s.SendChannels = sendChannels

	s.Rooms = make(map[connected_user.RoomName]*connected_user.Room)

	s.Rooms[connected_user.General] = connected_user.NewRoom(connected_user.General)
	s.Rooms[connected_user.Main] = connected_user.NewRoom(connected_user.Main)
	s.Rooms[connected_user.Fitness] = connected_user.NewRoom(connected_user.Fitness)
}

func (s *Server) BroadcastEvent(senderID int64, validEvent protocol.Event) {
	room := s.Rooms[validEvent.RouteRoomName]
	if room == nil {
		s.Logger.Debug("Room %v does not exist", validEvent.RouteRoomName)
		return
	}
	if room.ActiveUsers == nil {
		s.Logger.Debug("ActiveUsers is nil for room %v", validEvent.RouteRoomName)
		return
	}
	for userId := range room.ActiveUsers {
		if userId != senderID {
			s.RecvChannels[userId] <- validEvent
		}
	}
}

func (s *Server) ProcessUserState(u *connected_user.ConnectedUser, validEvent protocol.Event) (protocol.Event, error) {
	var roomName connected_user.RoomName
	var err error

	switch validEvent.Type {
	case protocol.EventTypeJoin:
		room := s.Rooms[roomName]
		err = u.JoinRoom(room)
		if err != nil {
			return protocol.Event{}, err
		}

	case protocol.EventTypeChat:
		if !u.InRoom {
			return protocol.Event{}, fmt.Errorf("User not in room, so can't send a message")
		}
		roomName = u.Room.RoomName
	case protocol.EventTypeLeave:
		roomName = u.Room.RoomName
		err = u.LeaveRoom()
		if err != nil {
			return protocol.Event{}, err
		}
	default:
		return protocol.Event{}, fmt.Errorf("The other message types are not supported")
	}

	validEvent.RouteRoomName = roomName
	return validEvent, nil
}

// we could to something even better send a ping through the connection to see if the user is still live
func HandleNetworkMessages(conn *net.TCPConn, networkMessages chan protocol.ClientMessage, serde *serde.Serde, logger *slog.Logger) {
	buf := make([]byte, 0, 200)
	defer close(networkMessages) // Close channel when goroutine exits
	for {
		msgBytes, err := network.ReadProtocol(conn, buf)
		if err != nil {
			logger.Info("Error reading from connection: %s\n", err)
			return
		}
		// Handle timeout/no data case
		if msgBytes == nil {
			continue // No message available, keep trying
		}

		if len(msgBytes) == 0 {
			continue // Empty message, skip
		}

		var msg protocol.ClientMessage
		err = serde.DecodeEncrypted(msgBytes, &msg)
		if err != nil {
			return
		}
		logger.Debug("Received network message", "msg", msg)
		networkMessages <- msg
	}
}

func (s *Server) HandleConn(conn *net.TCPConn, user *connected_user.ConnectedUser) {
	networkMessages := make(chan protocol.ClientMessage, 10)
	validEvents := make(chan protocol.Event, 10)
	errChan := make(chan protocol.ServerResponse)
	ackChan := make(chan protocol.ServerResponse)

	// handle input comming from the network
	go HandleNetworkMessages(conn, networkMessages, s.Serde, s.Logger)

	// does event validation
	go func() {
		for {
			cmsg := <-networkMessages

			var emptyRoom connected_user.RoomName
			event := protocol.EventFromMessage(cmsg, nil, emptyRoom)

			err := event.ContentValidation()
			if err != nil {
				errChan <- protocol.NewErrResponse(event, err)
			}

			s.Logger.Debug("Validated a new event", "event", event, "UserId", user.ID, "Username", user.Username)

			event, err = s.ProcessUserState(user, event)
			if err != nil {
				errChan <- protocol.NewErrResponse(event, err)
			}
			s.Logger.Debug("Validated event has valid user state", "event", event, "UserId", user.ID, "Username", user.Username)

			ackChan <- protocol.NewACKResponse(event)

			validEvents <- event
		}
	}()

	for {
		select {
		case event := <-validEvents:
			senderId := user.ID
			s.BroadcastEvent(senderId, event)
		case event := <-s.RecvChannels[user.ID]:

			response := event.ToResponse()
			data, err := s.Serde.EncodeEncrypted(&response)
			if err != nil {
				return
			}
			s.Logger.Debug("Writting a valid msg to conn:", "response", response, "UserID", user.ID, "username", user.Username)
			err = network.WriteProtocol(conn, data)

			if err != nil {
				s.Logger.Error("Failed to write data to conn", "err", err)
			}
		case errMsg := <-errChan:
			data, err := s.Serde.EncodeEncrypted(&errMsg)
			if err != nil {
				return
			}
			s.Logger.Debug("Writting an error msg to conn:", "msg", errMsg, "UserID", user.ID, "username", user.Username)
			err = network.WriteProtocol(conn, data)
			if err != nil {
				s.Logger.Error("Failed to write data to conn",
					"err", err,
					"remoteAddr", conn.RemoteAddr(),
				)
			}
		case ackMsg := <-ackChan:
			data, err := s.Serde.EncodeEncrypted(&ackMsg)
			if err != nil {
				return
			}
			s.Logger.Debug("Writting an error msg to conn:", "msg", ackMsg, "UserID", user.ID, "username", user.Username)
			err = network.WriteProtocol(conn, data)
			if err != nil {
				s.Logger.Error("Failed to write data to conn",
					"err", err,
					"remoteAddr", conn.RemoteAddr(),
				)
			}
		}
	}
}
