package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

// In this modules we should only can about the UI we can assume we are dealing with Messages and User structs

import (
	"fmt"
	"local-chat/internal/client"
	"local-chat/internal/logger"
	"local-chat/internal/message"
	"local-chat/internal/user"
	"log"
	"log/slog"
	"strings"

	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// At the top of your file

const gap = "\n\n"

func main() {
	p := tea.NewProgram(initialModel())

	// In main() or initialModel()
	debugFile, err := os.OpenFile(
		"~/code/local-chat/ui_debug.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err == nil {
		log.SetOutput(debugFile)
	}

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	client      *client.ChatClient
	user        *user.User
	logger      *slog.Logger
	connected   bool
	err         error
}

func initialModel() model {
	logger, err := logger.NewFileLogger(logger.UI_LOG_PATH)
	if err != nil {
		panic(err)
	}

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

	// Init the Chatclient

	client, user, err := client.InitUserSession(logger)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to init user session")
	}
	fmt.Printf("Did we get a user:")

	return model{
		textarea:    ta,
		messages:    []string{},
		viewport:    vp,
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
		client:      client,
		user:        user,
		logger:      logger,
	}
}

type ServerMsg message.Message

func ReadServerMsg(chat *client.ChatClient) tea.Msg {
	msg := <-chat.Incoming
	return ServerMsg(msg)
}
func listenForServerMsg(client *client.ChatClient) tea.Cmd { // Match your type
	return func() tea.Msg {
		return ReadServerMsg(client) // Use your existing function
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		listenForServerMsg(m.client), // Now returns tea.Cmd
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			msgType, err := ParseMsgType(msgTexted)
			if err != nil {
				//fmt.Println("Failed parsing the the msg type", err.Error())
			}

			message := MakeMsgFromType(m.user.Username, msgType, msgTexted)
			m.client.Outgoing <- message

			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	case ServerMsg:
		switch msg.Type {
		case message.Join, message.Leave, message.Error:
			m.messages = append(m.messages, m.senderStyle.Render(msg.Content))
		case message.Text:
			m.messages = append(m.messages, m.senderStyle.Render(fmt.Sprintf("%s: ", msg.Username))+msg.Content)
		}
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
		return m, tea.Batch(listenForIncomingMsg(m.client), tiCmd, vpCmd)

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
