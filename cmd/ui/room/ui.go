package ui

// A simple program demonstrating the text area component from the Bubbles
// component library.

// In this modules we should only can about the UI we can assume we are dealing with Messages and User structs

import (
	"errors"
	"fmt"
	"local-chat/cmd/ui/shared"
	"local-chat/cmd/ui/state"
	"local-chat/internal/client"
	"local-chat/internal/protocol"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// At the top of your file

const gap = "\n\n"

type (
	errMsg error
)

type RoomModel struct {
	viewport viewport.Model
	messages []string
	textarea textarea.Model

	senderStyle   lipgloss.Style
	receiverStyle lipgloss.Style
	joinedStyle   lipgloss.Style
	leftStyle     lipgloss.Style
	errStyle      lipgloss.Style

	client *client.ChatClient

	username string
	roomName string

	err error

	State *state.State
}

func (m RoomModel) SetUsername(username string) {
	m.username = username
}

func InitialRoomModel(username string, client *client.ChatClient, state *state.State) RoomModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return RoomModel{
		textarea: ta,
		messages: []string{},
		viewport: vp,

		senderStyle:   senderStyle,
		receiverStyle: receiverStyle,
		joinedStyle:   joinedStyle,
		leftStyle:     leftStyle,
		errStyle:      errStyle,

		err: nil,

		username: username,
		client:   client,
		State:    state,
	}
}

func (m RoomModel) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		shared.ListenForServerMsg(m.client),
	)
}

func (m RoomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			msgTexted := m.textarea.Value()
			msgType, roomName, err := ParseMsgType(msgTexted)
			if err != nil {
				m.err = err
				return m, nil
			}
			message := BuildClientMessage(msgType, m.username, msgTexted, roomName)

			m.client.Outgoing <- message

			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	case shared.ServerACK:
		switch msg.ActionType {
		case protocol.EventTypeJoin:
		//	we need to switch to the room state
		case protocol.EventTypeLeave:
			// we need to go to the dash screen
			*m.State = state.Dash
		case protocol.EventTypeChat:
			message := m.senderStyle.Render(fmt.Sprintf("%s: ", msg.Username)) + msg.Content
			m.messages = append(m.messages, message)

			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.viewport.GotoBottom()

		}
	case shared.ServerErr:
		message := m.errStyle.Render(msg.Content)
		m.err = errors.New(message)

	case shared.ServerBrodcast:
		var message string
		switch msg.ActionType {
		case protocol.EventTypeChatBroadcast:
			message = m.receiverStyle.Render(fmt.Sprintf("%s: ", msg.Username)) + msg.Content
		case protocol.EventTypeUserJoined:
			message = m.joinedStyle.Render(fmt.Sprintf("%s: ", msg.Username)) + msg.Content
		case protocol.EventTypeUserLeft:
			message = m.leftStyle.Render(fmt.Sprintf("%s: ", msg.Username)) + msg.Content
		default:
			panic(fmt.Errorf("Wrong message brodcast received: %v", msg.Type))
		}
		m.messages = append(m.messages, message)

	//switch msg.Type {
	//case message.Join, message.Leave, message.Error:
	//	m.messages = append(m.messages, m.senderStyle.Render(msg.Content))
	//case message.Text:
	//	m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%s: ", msg.Username))+msg.Content)
	//}
	//m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
	//m.viewport.GotoBottom()
	//return m, tea.Batch(listenForIncomingMsg(m.client), tiCmd, vpCmd)

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(shared.ListenForServerMsg(m.client), tiCmd, vpCmd)
}

func (m RoomModel) View() string {
	// If there's an error, show it prominently
	if m.err != nil {
		errorView := m.errStyle.Render(fmt.Sprintf("Error: %s", m.err.Error()))
		m.err = nil // Reset the error after displaying it
		return fmt.Sprintf(
			"%s%s%s%s%s",
			m.viewport.View(),
			gap,
			errorView,
			gap,
			m.textarea.View(),
		)
	}

	// Normal view when no error
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
