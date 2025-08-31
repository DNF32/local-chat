package shared

import (
	"local-chat/internal/client"
	"local-chat/internal/protocol"

	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type ServerACK protocol.ServerResponse
type ServerErr protocol.ServerResponse
type ServerLogin protocol.ServerResponse
type ServerBrodcast protocol.ServerResponse

func ReadServerMsg(chat *client.ChatClient) tea.Msg {
	msg := <-chat.Incoming
	switch msg.Type {
	case protocol.EventTypeAck:
		return ServerACK(msg)
	case protocol.EventTypeError:
		return ServerErr(msg)
	case protocol.EventTypeChatBroadcast, protocol.EventTypeUserLeft, protocol.EventTypeUserJoined:
		return ServerBrodcast(msg)
	case protocol.EventTypeSucessLogin, protocol.EventTypeFailedLogin:
		return ServerLogin(msg)
	default:
		panic(fmt.Errorf("Server sent a ServerResponse with wrong type %s", msg.Type))
	}
}

func ListenForServerMsg(client *client.ChatClient) tea.Cmd { // Match your type
	return func() tea.Msg {
		return ReadServerMsg(client) // Use your existing function
	}
}
