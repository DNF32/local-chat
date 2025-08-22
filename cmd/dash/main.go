package main

import (
	"fmt"
	"local-chat/internal/logger"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	//"github.com/charmbracelet/lipgloss"
)

const dash = `$$\                                   $$\        $$$$$$\  $$\                  $$\     $$\                         
$$ |                                  $$ |      $$  __$$\ $$ |                 $$ |    $$ |                        
$$ |      $$$$$$\   $$$$$$$\ $$$$$$\  $$ |      $$ /  \__|$$$$$$$\   $$$$$$\ $$$$$$\ $$$$$$\    $$$$$$\   $$$$$$\  
$$ |     $$  __$$\ $$  _____|\____$$\ $$ |      $$ |      $$  __$$\  \____$$\\_$$  _|\_$$  _|  $$  __$$\ $$  __$$\ 
$$ |     $$ /  $$ |$$ /      $$$$$$$ |$$ |      $$ |      $$ |  $$ | $$$$$$$ | $$ |    $$ |    $$$$$$$$ |$$ |  \__|
$$ |     $$ |  $$ |$$ |     $$  __$$ |$$ |      $$ |  $$\ $$ |  $$ |$$  __$$ | $$ |$$\ $$ |$$\ $$   ____|$$ |      
$$$$$$$$\\$$$$$$  |\$$$$$$$\\$$$$$$$ |$$ |      \$$$$$$  |$$ |  $$ |\$$$$$$$ | \$$$$  |\$$$$  |\$$$$$$$\ $$ |      
\________|\______/  \_______|\_______|\__|       \______/ \__|  \__| \_______|  \____/  \____/  \_______|\__|      
                                                                                                                   
                                                                                                                   `

type Dash struct {
	textarea textarea.Model
	width    int
	height   int
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

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling

	ta.ShowLineNumbers = false

	ta.KeyMap.InsertNewline.SetEnabled(false)
	return Dash{textarea: ta,
		width:  80,
		height: 24}
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

func (m Dash) View() string {
	textarea := m.textarea.View()
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			dash,
			"",
			textarea,
		),
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	p.Run()
}
