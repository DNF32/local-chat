package dash

import (
	"errors"
	"fmt"
	"local-chat/cmd/ui/shared"
	"local-chat/cmd/ui/state"
	"local-chat/internal/client"
	"local-chat/internal/protocol"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const dash2 = ` _       ___     __   ____  _             __  __ __   ____  ______  ______    ___  ____  
| |     /   \   /  ] /    || |           /  ]|  |  | /    ||      ||      |  /  _]|    \ 
| |    |     | /  / |  o  || |          /  / |  |  ||  o  ||      ||      | /  [_ |  D  )
| |___ |  O  |/  /  |     || |___      /  /  |  _  ||     ||_|  |_||_|  |_||    _]|    / 
|     ||     /   \_ |  _  ||     |    /   \_ |  |  ||  _  |  |  |    |  |  |   [_ |    \ 
|     ||     \     ||  |  ||     |    \     ||  |  ||  |  |  |  |    |  |  |     ||  .  \
|_____| \___/ \____||__|__||_____|     \____||__|__||__|__|  |__|    |__|  |_____||__|\_|
                                                                                         `

const dash = ` ______                   _______           _______ _________
(  __  \ |\     /|       (  ____ \|\     /|(  ___  )\__   __/
| (  \  )| )   ( |       | (    \/| )   ( || (   ) |   ) (   
| |   ) || |   | | _____ | |      | (___) || (___) |   | |   
| |   | || |   | |(_____)| |      |  ___  ||  ___  |   | |   
| |   ) || |   | |       | |      | (   ) || (   ) |   | |   
| (__/  )| (___) |       | (____/\| )   ( || )   ( |   | |   
(______/ (_______)       (_______/|/     \||/     \|   )_(   
                                                             `

type User struct {
	Username string
	InRoom   bool
	RoomName string
}

type Dash struct {
	textarea textarea.Model
	width    int
	height   int
	errStyle lipgloss.Style

	username string
	client   *client.ChatClient

	err error
}

func (m *Dash) SetUsername(username string) {
	m.username = username
}

func (m Dash) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		shared.ListenForServerMsg(m.client), // Now returns tea.Cmd
	)
}

func InitialModel(username string, client *client.ChatClient, width, height int) Dash {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "\uf18e :  "
	ta.CharLimit = 280

	ta.SetWidth(20)
	ta.SetHeight(1)

	// Remove cursor line styling

	ta.ShowLineNumbers = false

	ta.KeyMap.InsertNewline.SetEnabled(false)
	return Dash{textarea: ta,
		width:    width,
		height:   height,
		client:   client,
		username: username,

		errStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Background(lipgloss.Color("#330000")).
			Bold(true),
	}
}

func (m Dash) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	// clean err if it's present
	if m.err != nil {
		m.err = nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(msg.Width)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			msgTexted := m.textarea.Value()
			msgType, roomName, err := shared.ParseMsgType(msgTexted)
			if err != nil {
				m.err = err
				return m, nil
			}
			message := shared.BuildClientMessage(msgType, m.username, msgTexted, roomName)

			m.client.Outgoing <- message
			m.textarea.Reset()

		}
	case shared.ServerACK:
		switch msg.ActionType {
		case protocol.EventTypeJoin:
			stateChangeCmd := func() tea.Msg {
				return state.StateChangeMsg{NewState: state.Room, Username: string(msg.Username)}
			}
			return m, stateChangeCmd
		case protocol.EventTypeChat:
			panic("Dash received a NetworkACk of type EventTypeChat, this is an invalida state")
		}
	case shared.ServerErr:
		message := m.errStyle.Render(msg.Content)
		m.err = errors.New(message)
	}
	return m, tea.Batch(tiCmd, shared.ListenForServerMsg(m.client))
}

func (m Dash) View() string {
	textarea := m.textarea.View()
	// Style the textarea to make it smaller and more like a text box
	textBox := lipgloss.NewStyle().
		Width(int(float64(m.width)*0.5)). // 50% of window width
		Height(3).                        // Set a small height (like 3 lines)
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 1). // Small horizontal padding inside the box
		Render(textarea)

	// Simplified user info
	username := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")). // Blue
		Render(m.username)

	// Available channels
	channels := lipgloss.NewStyle().
		Foreground(lipgloss.Color("246")). // Gray
		Render("Channels: Main, General, Fitness")

	userInfo := lipgloss.JoinVertical(
		lipgloss.Center,
		username,
		channels,
	)

	userBox := lipgloss.NewStyle().
		Width(40).
		Height(3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1).
		Align(lipgloss.Center).
		Render(userInfo)

	userArea := lipgloss.JoinHorizontal(lipgloss.Center, textBox, userBox)

	// Build the main content
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		dash,
		"\n\n\n\n",
		userArea,
	)

	// Add error message if present
	if m.err != nil {
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			content,
			"\n",
			m.err.Error(),
		)

		// Clear the error after displaying
		m.err = nil
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}
