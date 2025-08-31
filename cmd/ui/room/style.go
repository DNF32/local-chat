package ui

import "github.com/charmbracelet/lipgloss"

var (
	// sender messages (you)
	senderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00CFFF")). // bright electric blue
			Bold(true)

	// received messages (others)
	receiverStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B084F7")) // soft purple to contrast

	// user joined
	joinedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF9F")). // neon aqua/green
			Background(lipgloss.Color("#002233")).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FF9F")).
			Bold(true)

	// user left
	leftStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4B91")). // hot pink / magenta
			Background(lipgloss.Color("#220011")).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF4B91")).
			Italic(true)

	// error
	errStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Background(lipgloss.Color("#330000")).
			Bold(true)
)
