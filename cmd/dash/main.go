package main

import (
	"fmt"
	"local-chat/internal/logger"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	//"github.com/charmbracelet/lipgloss"
)

const dash = ` _       ___     __   ____  _             __  __ __   ____  ______  ______    ___  ____  
| |     /   \   /  ] /    || |           /  ]|  |  | /    ||      ||      |  /  _]|    \ 
| |    |     | /  / |  o  || |          /  / |  |  ||  o  ||      ||      | /  [_ |  D  )
| |___ |  O  |/  /  |     || |___      /  /  |  _  ||     ||_|  |_||_|  |_||    _]|    / 
|     ||     /   \_ |  _  ||     |    /   \_ |  |  ||  _  |  |  |    |  |  |   [_ |    \ 
|     ||     \     ||  |  ||     |    \     ||  |  ||  |  |  |  |    |  |  |     ||  .  \
|_____| \___/ \____||__|__||_____|     \____||__|__||__|__|  |__|    |__|  |_____||__|\_|
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
	user     User
}

func (m Dash) Init() tea.Cmd {
	return textarea.Blink
}

func initialModel() Dash {
	_, err := logger.NewFileLogger(logger.UI_LOG_PATH)
	if err != nil {
		panic(err)
	}

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
		width:  80,
		height: 24,
		user: User{Username: "Tester",
			InRoom:   false,
			RoomName: ""},
	}
}

func (m Dash) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		tiCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)

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
			_ = m.textarea.Value()
			m.textarea.Reset()

			return m, tiCmd
		}
	}
	return m, tiCmd
}

func (u User) View() string {
	username := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")). // Blue
		Render(u.Username)

	var roomStatus string
	if u.InRoom {
		roomStatus = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")). // Green
			Render(fmt.Sprintf("🟢 %s", u.RoomName))
	} else {
		roomStatus = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Render("🔴 Not connected")
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		username,
		"",
		roomStatus,
	)
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

	userBox := lipgloss.NewStyle().
		Width(40). // 50% of window width
		Height(3). // Set a small height (like 3 lines)
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1). // Small horizontal padding inside the box
		Align(lipgloss.Center).
		Render(m.user.View())

	userArea := lipgloss.JoinHorizontal(lipgloss.Center, textBox, userBox)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			dash,
			"\n\n\n\n",
			userArea,
		),
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	p.Run()
}
