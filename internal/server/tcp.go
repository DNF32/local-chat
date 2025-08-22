package server

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"local-chat/internal/logger"
	"local-chat/internal/message"
	"local-chat/internal/network"
	"local-chat/internal/user"
)

type Event struct {
	Type      message.MessageType `json:"type"`
	Username  string              `json:"username"`
	Action    string
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	RoomName  user.RoomName
}

func (e *Event) ToMsg() message.Message {
	var content string
	switch e.Type {
	case message.Text, message.Undefined, message.Error:
		content = e.Content
	case message.Join, message.Leave:
		content = fmt.Sprintf("%s %s Room %s", e.Username, e.Action, e.RoomName)
	}
	return message.Message{Type: e.Type,
		Username:  e.Username,
		Content:   content,
		Timestamp: e.Timestamp}
}

func (e *Event) EqualTo(other *Event) (bool, string) {
	if e.Type != other.Type {
		return false, "Type"
	}
	if e.Username != other.Username {
		return false, "Username"
	}
	if e.Action != other.Action {
		return false, "Action"
	}
	if e.Content != other.Content {
		return false, "Content"
	}
	if e.RoomName != other.RoomName {
		return false, "RoomName"
	}
	return true, ""
}

func ValidateNewEvent(m message.Message) (Event, error) {
	msgType, content := m.ParseContent()
	if msgType == message.Undefined {
		return Event{}, fmt.Errorf("Failed to parse message content into server Event")
	}
	return Event{
		Type:     m.Type,
		Username: m.Username,
		Action:   string(m.Type),
		Content:  content, Timestamp: m.Timestamp}, nil
}

var id int = 1

type Server struct {
	Users        map[int]*user.User
	Rooms        map[user.RoomName]*user.Room
	SendChannels map[int]chan Event // for sending TO users
	RecvChannels map[int]chan Event // for receiving FROM users

	Listener *net.TCPListener
	Logger   *slog.Logger
}

// This function produces a new event data
// Modifies user state according to the type of event

func (s *Server) BroadcastEvent(validEvent Event) {
	room := s.Rooms[validEvent.RoomName]
	if room == nil {
		s.Logger.Debug("Room %v does not exist", validEvent.RoomName)
		return
	}
	if room.ActiveUsers == nil {
		s.Logger.Debug("ActiveUsers is nil for room %v", validEvent.RoomName)
		return
	}
	for userId := range room.ActiveUsers {
		s.RecvChannels[userId] <- validEvent
	}
}

func (s *Server) ProcessUserEvent(u *user.User, validEvent Event) (Event, error) {
	var roomName user.RoomName
	var err error

	switch validEvent.Type {
	case message.Join:
		roomName, err = user.ValidateRoom(validEvent.Content)
		if err != nil {
			return Event{}, err
		}
		room := s.Rooms[roomName]
		err = u.JoinRoom(room)
		if err != nil {
			return Event{}, err
		}

	case message.Text:
		if !u.InRoom {
			return Event{}, fmt.Errorf("User not in room, so can't send a message")
		}
		roomName = u.Room.RoomName
	case message.Leave:
		roomName, err = u.LeaveRoom()
		if err != nil {
			return Event{}, err
		}
	default:
		return Event{}, fmt.Errorf("The other message types are not supported")
	}

	validEvent.RoomName = roomName
	return validEvent, nil
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

	sendChannels := make(map[int]chan Event)
	recvChannels := make(map[int]chan Event)

	s.RecvChannels = recvChannels
	s.SendChannels = sendChannels

	s.Rooms = make(map[user.RoomName]*user.Room)

	s.Rooms[user.General] = user.NewRoom(user.General)
	s.Rooms[user.Main] = user.NewRoom(user.Main)
	s.Rooms[user.Fitness] = user.NewRoom(user.Fitness)
}

// This enables the command list active users
// we could to something even better send a ping through the connection to see if the user is still live

// INFO: This is the entry point of the messages
func HandleNetworkMessages(conn *net.TCPConn, networkMessages chan message.Message, logger *slog.Logger) {
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
		var msg message.Message
		err = msg.Decode(msgBytes)
		if err != nil {
			return
		}
		logger.Debug("Received network message", "msg", msg)
		networkMessages <- msg
	}
}

func (s *Server) HandleConn(conn *net.TCPConn, user *user.User) {
	networkMessages := make(chan message.Message, 10)
	validEvents := make(chan Event, 10)
	errChan := make(chan message.Message)

	// handle input comming from the network
	go HandleNetworkMessages(conn, networkMessages, s.Logger)

	// does event validation
	go func() {
		for {
			msg := <-networkMessages
			event, err := ValidateNewEvent(msg)
			if err != nil {
				errChan <- message.NewErr(user.Username, err)
			}

			s.Logger.Debug("Validated a new event", "event", event, "UserId", user.ID, "Username", user.Username)

			event, err = s.ProcessUserEvent(user, event)
			if err != nil {
				errChan <- message.NewErr(user.Username, err)
			}
			s.Logger.Debug("Validated event has valid user state", "event", event, "UserId", user.ID, "Username", user.Username)

			validEvents <- event
		}
	}()

	for {
		select {
		case event := <-validEvents:
			s.BroadcastEvent(event)
		case event := <-s.RecvChannels[user.ID]:
			msg := event.ToMsg()
			data, err := msg.Encode()
			if err != nil {
				return
			}
			s.Logger.Debug("Writting a valid msg to conn:", "msg", msg, "UserID", user.ID, "username", user.Username)
			for len(data) > 0 {
				n, err := conn.Write(data)
				if err != nil {
					return
				}
				data = data[n:]
			}
		case errMsg := <-errChan:
			data, err := errMsg.Encode()
			if err != nil {
				return
			}
			s.Logger.Debug("Writting an error msg to conn:", "msg", errMsg, "UserID", user.ID, "username", user.Username)
			for len(data) > 0 {
				n, err := conn.Write(data)
				if err != nil {
					return
				}
				data = data[n:]
			}
		}
	}
}

func main() {
	var id int = 1

	s := Server{}
	s.Start()
	for {
		conn, err := s.Listener.AcceptTCP()
		if err != nil {
			return
		}

		user, err := InitUser(conn, id)
		if err != nil {
			// Failed to init this user continue in the loop to handle a new connection
			continue
		}
		id++

		userSend := make(chan Event, 20)
		userRecv := make(chan Event, 20)
		s.SendChannels[user.ID] = userSend
		s.RecvChannels[user.ID] = userRecv
		fmt.Printf("Created send/recv channels for user %d\n", user.ID)

		go s.HandleConn(conn, user)
		fmt.Printf("Started goroutine to handle connection for user %d\n", user.ID)
	}
}

func InitUser(conn *net.TCPConn, id int) (*user.User, error) {
	initMsgBuf := make([]byte, 0, 200)

	// Read bytes send
	data, err := network.ReadProtocol(conn, initMsgBuf)
	if err != nil {
		fmt.Println("Error reading init message")
		return nil, fmt.Errorf("failed to read init message: %w", err)
	}
	fmt.Printf("Read %d bytes for init message\n", len(data))

	// Decode the message sent
	var initMsg message.Message
	err = initMsg.Decode(data)
	if err != nil {
		fmt.Printf("Failed to decode init message: %v\n", err)
		return nil, fmt.Errorf("failed to decode init message: %w", err)
	}
	fmt.Printf("Decoded init message: Username=%s\n", initMsg.Username)

	// Initing user and writting to the conn
	user := &user.User{
		ID:       id,
		Username: initMsg.Username,
	}

	data, _ = user.Encode()
	fmt.Printf("Sending user ID %d to client\n", user.ID)

	for len(data) > 0 {
		n, err := conn.Write(data)
		if err != nil {
			fmt.Printf("Error writing user ID to client: %v\n", err)
			return nil, fmt.Errorf("failed to send user data: %w", err)
		}
		data = data[n:]
	}
	return user, nil
}
