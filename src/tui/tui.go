package tui

import (
	"fmt"
	"os"

	"github.com/Ryan-Har/csnip/common/models"
	"github.com/Ryan-Har/csnip/database"
	"github.com/Ryan-Har/csnip/tui/editTui"
	"github.com/Ryan-Har/csnip/tui/viewTui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TUIOpts struct {
	Theme string
}

type state int

const (
	viewState state = iota
	addState
	editState
)

type model struct {
	db           database.DatabaseInteractions // used to interact with the database
	currentState state                         // view currently being used
	viewModel    tea.Model
	editModel    tea.Model
	Width        int
	Height       int
}

// Messages confirm users intent
type ClearScreenMsg struct{}

func initialModel(db database.DatabaseInteractions) model {
	return model{
		db:           db,
		currentState: viewState,
		viewModel:    viewTui.New(db),
		editModel:    editTui.New(db),
	}
}

func (m model) Init() tea.Cmd {
	return m.viewModel.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case viewTui.SingleEditMsg:
		m.editModel = editTui.New(m.db,
			editTui.WithCodeSnippet(models.CodeSnippet(msg)),
			editTui.InEditMode())
		m.currentState = editState
		return m, tea.Batch(m.clearScreen(), m.editModel.Init())
	case viewTui.SingleViewMsg:
		m.editModel = editTui.New(m.db,
			editTui.WithCodeSnippet(models.CodeSnippet(msg)))
		m.currentState = editState
		return m, tea.Batch(m.clearScreen(), m.editModel.Init())
	case viewTui.SingleAddMsg:
		m.editModel = editTui.New(m.db,
			editTui.InEditMode())
		m.currentState = editState
		return m, tea.Batch(m.clearScreen(), m.editModel.Init())
	case editTui.WindowSizeReqMsg:
		return m, m.sendWindowsSizeMessage()
	case editTui.ReturnToViewMsg:
		m.currentState = viewState
		return m, m.viewModel.Init()
	// propogate windowSizeMsg so subModels can size correctly
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		newViewModel, newViewCmd := m.viewModel.Update(msg)
		m.viewModel = newViewModel
		newEditModel, newEditCmd := m.editModel.Update(msg)
		m.editModel = newEditModel
		cmds = append(cmds, newViewCmd, newEditCmd)
	// always quit if requested
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case ClearScreenMsg:
		return m, tea.ClearScreen
	}

	switch m.currentState {
	case viewState:
		newViewModel, newCmd := m.viewModel.Update(msg)
		m.viewModel = newViewModel
		cmd = newCmd
	case editState:
		newEditModel, newCmd := m.editModel.Update(msg)
		m.editModel = newEditModel
		cmd = newCmd
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m model) View() string {
	switch m.currentState {
	case viewState:
		return m.viewModel.View()
	case editState:
		return m.editModel.View()
	default:
		return "Unknown state"
	}
}

func (t *TUIOpts) Run(db database.DatabaseInteractions) {
	p := tea.NewProgram(initialModel(db), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) clearScreen() tea.Cmd {
	return func() tea.Msg {
		return ClearScreenMsg{}
	}
}

func (m model) sendWindowsSizeMessage() tea.Cmd {
	return func() tea.Msg {
		msg := tea.WindowSizeMsg{
			Height: m.Height,
			Width:  m.Width,
		}
		return msg
	}
}

// func deleteSnippetFromDB(db database.DatabaseInteractions, id uuid.UUID) tea.Cmd {
// 	return func() tea.Msg {
// 		err := db.DeleteSnippetByUUID(id)
// 		if err != nil {
// 			return errMsg{err}
// 		}
// 		return
// 	}
// }

// var baseStyle = lipgloss.NewStyle().
// 	BorderStyle(lipgloss.NormalBorder()).
// 	BorderForeground(lipgloss.Color("240"))

// type KeyMap struct {
// 	Quit  key.Binding
// 	Up    key.Binding
// 	Down  key.Binding
// 	Enter key.Binding
// 	Esc   key.Binding
// }

// var defaultKeyMap = KeyMap{
// 	Quit: key.NewBinding(
// 		key.WithKeys("ctrl+c", "q"),
// 		key.WithHelp("exit", "quit"),
// 	),
// 	Up: key.NewBinding(
// 		key.WithKeys("k", "up"),
// 		key.WithHelp("↑/k", "move up"),
// 	),
// 	Down: key.NewBinding(
// 		key.WithKeys("j", "down"),
// 		key.WithHelp("↓/j", "move down"),
// 	),
// 	Enter: key.NewBinding(
// 		key.WithKeys("enter", "v"),
// 		key.WithHelp("Enter", "view"),
// 	),
// 	Esc: key.NewBinding(
// 		key.WithKeys("esc"),
// 		key.WithHelp("escape", "back"),
// 	),
// }
