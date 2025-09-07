package main

import (
	"local-chat/cmd/ui/dash"
	"local-chat/cmd/ui/login"
	"local-chat/cmd/ui/room"
	"local-chat/cmd/ui/state"
	"local-chat/internal/client"
	//"local-chat/internal/logger"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Router struct {
	State  state.State
	width  int
	height int

	room   ui.RoomModel
	login  login.LoginModel
	dash   dash.Dash
	client *client.ChatClient
}

func InitialModel() Router {
	state := state.Login

	client, _ := client.NewChatClient()
	go client.HandleInput()
	go client.HandleOutput()
	//logger, _ := logger.NewFileLogger(logger.UI_LOG_PATH)

	room := ui.InitialRoomModel("", client, 30, 5)
	login := login.InitialLogin(client, 30, 5)
	dash := dash.InitialModel("", client, 30, 5)
	return Router{
		State:  state,
		room:   room,
		login:  login,
		dash:   dash,
		client: client,
		width:  30,
		height: 5,
	}
}

func (r Router) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle messages globally first
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		modelLogin, c := r.login.Update(msg)
		r.login = modelLogin.(login.LoginModel)
		cmdLogin := c

		modelDash, c := r.dash.Update(msg)
		r.dash = modelDash.(dash.Dash)
		cmdDash := c

		modelRoom, c := r.room.Update(msg)
		r.room = modelRoom.(ui.RoomModel)
		cmdRoom := c

		return r, tea.Batch(cmdLogin, cmdDash, cmdRoom)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return r, tea.Quit
		}
	case state.StateChangeMsg:
		// Handle state changes globally FIRST
		r.State = msg.NewState

		switch msg.NewState {
		case state.Dash:
			if msg.Username != "" {
				r.dash = dash.InitialModel(msg.Username, r.client, r.width, r.height)
				// Return immediately with the dash init command
				return r, r.dash.Init()
			}
		case state.Room:
			// Handle room initialization if needed
			r.room = ui.InitialRoomModel(msg.Username, r.client, r.width, r.height)
			return r, r.room.Init()
		case state.Login:
			// Handle login initialization if needed
			r.login = login.InitialLogin(r.client, r.width, r.height)
			return r, r.login.Init()
		}
		// For any state change, return immediately to trigger re-render
		return r, nil
	}

	// Route messages to current state AFTER handling state changes
	switch r.State {
	case state.Login:
		model, c := r.login.Update(msg)
		r.login = model.(login.LoginModel)
		cmd = c
	case state.Dash:
		model, c := r.dash.Update(msg)
		r.dash = model.(dash.Dash)
		cmd = c
	case state.Room:
		model, c := r.room.Update(msg)
		r.room = model.(ui.RoomModel)
		cmd = c
	}

	return r, cmd
}
func (r Router) View() string {
	// Add debug logging with proper structured format
	stateStr := ""
	switch r.State {
	case state.Login:
		stateStr = "Login"
	case state.Dash:
		stateStr = "Dash"
	case state.Room:
		stateStr = "Room"
	default:
		stateStr = "Unknown"
	}

	r.client.Logger.Info("Rendering view", "state", stateStr, "stateValue", int(r.State))

	// Route to appropriate view based on current state
	switch r.State {
	case state.Login:
		r.client.Logger.Info("Rendering login view")
		return r.login.View()
	case state.Dash:
		r.client.Logger.Info("Rendering dash view")
		return r.dash.View()
	case state.Room:
		r.client.Logger.Info("Rendering room view")
		return r.room.View()
	default:
		return ""
	}
}

func (r Router) Init() tea.Cmd {
	var cmds []tea.Cmd

	// Initialize all models
	cmds = append(cmds, r.login.Init())
	cmds = append(cmds, r.dash.Init())
	cmds = append(cmds, r.room.Init())

	// Return batched commands
	return tea.Batch(cmds...)
}

func main() {
	p := tea.NewProgram(InitialModel())

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
