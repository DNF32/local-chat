package login

import (
	"fmt"
	"local-chat/cmd/ui/shared"
	"local-chat/cmd/ui/state"
	"local-chat/internal/client"
	"local-chat/internal/protocol"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoginState struct {
	sucess   bool
	username client.Username
}

func LoginCMD(m LoginModel, username, password string) tea.Cmd {
	return func() tea.Msg {
		m.client.SendAuthRequest(username, password)
		return nil
	}
}

type (
	errMsg error
)

const (
	us = iota
	pass
)

const (
	eletricBlue = lipgloss.Color("#00D4FF")
	darkGray    = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(eletricBlue)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type LoginModel struct {
	inputs  []textinput.Model
	focused int
	width   int
	height  int
	err     error

	client *client.ChatClient

	State *state.State
}

func InitialLogin(client *client.ChatClient, state *state.State) LoginModel {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[us] = textinput.New()
	inputs[us].Placeholder = "User1"
	inputs[us].Focus()
	inputs[us].CharLimit = 20
	inputs[us].Width = 30
	inputs[us].Prompt = ""

	inputs[pass] = textinput.New()
	inputs[pass].Placeholder = "******"
	inputs[pass].CharLimit = 15
	inputs[pass].Width = 15
	inputs[pass].Prompt = ""
	inputs[pass].EchoMode = textinput.EchoPassword

	return LoginModel{
		inputs:  inputs,
		width:   80,
		height:  24,
		focused: 0,
		err:     nil,
		State:   state,
		client:  client,
	}
}

func (m LoginModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, shared.ListenForServerMsg(m.client))
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	// Add debug logging for all messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				m.client.Logger.Info("Submitting login form")
				return m, LoginCMD(m, m.inputs[us].Value(), m.inputs[pass].Value())
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	case LoginState:
		m.client.Logger.Info("Received LoginState message", "success", msg.sucess, "username", string(msg.username))
		if !msg.sucess {
			m.client.Logger.Info("Failed Login")
			for i := range m.inputs {
				m.inputs[i].SetValue("")
				m.focused = 0
				m.inputs[m.focused].Focus()
			}
		} else {
			m.client.Logger.Info("Login successful, sending state change message")
			for i := range m.inputs {
				m.inputs[i].SetValue("")
				m.focused = 0
				m.inputs[m.focused].Focus()
			}
			// Create the state change command
			stateChangeCmd := func() tea.Msg {
				return state.StateChangeMsg{NewState: state.Dash, Username: string(msg.username)}
			}
			m.client.Logger.Info("Returning state change command")
			return m, stateChangeCmd
		}
	case shared.ServerLogin:
		switch msg.Type {
		case protocol.EventTypeFailedLogin:
			m.client.Logger.Info("Failed Login")
			for i := range m.inputs {
				m.inputs[i].SetValue("")
				m.focused = 0
				m.inputs[m.focused].Focus()
			}
		case protocol.EventTypeSucessLogin:
			m.client.Logger.Info("Login successful, sending state change message")
			for i := range m.inputs {
				m.inputs[i].SetValue("")
				m.focused = 0
				m.inputs[m.focused].Focus()
			}
			// Create the state change command
			stateChangeCmd := func() tea.Msg {
				return state.StateChangeMsg{NewState: state.Dash, Username: msg.Username}
			}
			m.client.Logger.Info("Returning state change command")
			return m, stateChangeCmd
		}

	case errMsg:
		m.client.Logger.Info("Received error message", "error", msg)
		m.err = msg
		return m, nil
	default:
		// Log other message types for debugging
		m.client.Logger.Info("Received other message type", "type", fmt.Sprintf("%T", msg))
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m LoginModel) View() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			dash,
			"\n\n\n\n",
			m.BuildString(),
		),
	)

}

func (m LoginModel) BuildString() string {
	return lipgloss.JoinVertical(
		lipgloss.Left, // keep label/input aligned left relative to each other
		inputStyle.Width(m.inputs[us].Width).Render("Username"),
		m.inputs[us].View(),
		"\n",
		inputStyle.Width(m.inputs[pass].Width).Render("Password"),
		m.inputs[pass].View(),
	)
}

// nextInput focuses the next input field
func (m *LoginModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *LoginModel) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
