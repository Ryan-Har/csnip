package viewTui

import (
	"errors"
	"fmt"

	"github.com/Ryan-Har/csnip/common/models"
	"github.com/Ryan-Har/csnip/database"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	CodeSnippets      []models.CodeSnippet          // code snippets currently being displayed
	CodeSnippetsTable table.Model                   // table representation of the code snippets
	FilterLang        string                        // current language filter
	FilterTag         string                        // current filter tag
	db                database.DatabaseInteractions // used to interact with the database
	keys              KeyMap                        // KeyMap holding available keys
	help              help.Model
	height            int // height of the terminal window
	width             int // width of the terminal window
}

// Messages for use within this model
type FetchSnippetsReqMsg struct{}
type SnippetsFetchedMsg []models.CodeSnippet
type SnippetsTableMsg table.Model
type ErrMsg struct {
	Err error
}

// Messages for use in the main model
type SingleViewMsg models.CodeSnippet
type SingleEditMsg models.CodeSnippet

// creates new tea.Model interface
func New(db database.DatabaseInteractions) tea.Model {
	return Model{
		CodeSnippets:      []models.CodeSnippet{},
		CodeSnippetsTable: table.Model{},
		FilterLang:        "",
		FilterTag:         "",
		db:                db,
		keys:              keyMap,
		help:              help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.getCodeSnippet()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrMsg:
		panic(msg.Err.Error())
	case FetchSnippetsReqMsg:
		return m, m.getCodeSnippet()
	case SnippetsFetchedMsg:
		m.CodeSnippets = msg
		return m, m.generateSnippetsTable(msg)
	case SnippetsTableMsg:
		m.CodeSnippetsTable = table.Model(msg)
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, m.generateSnippetsTable(m.CodeSnippets)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Down):
			m.CodeSnippetsTable.MoveDown(1)
		case key.Matches(msg, m.keys.Up):
			m.CodeSnippetsTable.MoveUp(1)
		case key.Matches(msg, m.keys.View):
			focused := m.CodeSnippets[m.CodeSnippetsTable.Cursor()]
			return m, m.openForViewing(focused)
		case key.Matches(msg, m.keys.Edit):
			focused := m.CodeSnippets[m.CodeSnippetsTable.Cursor()]
			return m, m.openForEditing(focused)
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}

		// case "d", "delete":
		// 	m.focusedSnippet = m.codeSnippets[m.codeSnippetsTable.Cursor()]
		// 	m.currentPage = confirmationPage
		// 	m.confirmDelete = true
		// 	return m, nil

	}
	return m, nil
}

func (m Model) View() string {
	var tableMsg string
	if len(m.CodeSnippets) == 0 {
		tableMsg = "No code Snippets Found\n"
	} else {
		tableMsg = baseStyle.Render(m.CodeSnippetsTable.View()) + "\n"
	}
	helpView := m.help.View(m.keys)

	return tableMsg + helpView
}

type KeyMap struct {
	Quit key.Binding
	Up   key.Binding
	Down key.Binding
	View key.Binding
	Edit key.Binding
	Add  key.Binding
	Help key.Binding
}

var keyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	View: key.NewBinding(
		key.WithKeys("enter", "v"),
		key.WithHelp("Enter/v", "view"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.View, k.Edit},
		{k.Help, k.Quit},
	}
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m Model) getCodeSnippet() tea.Cmd {
	return func() tea.Msg {
		if m.FilterLang != "" && m.FilterTag != "" {
			snips, err := m.db.GetSnippetsByLanguageAndTag(m.FilterLang, m.FilterTag)
			if err != nil {
				if errors.Is(err, database.ErrNoSnippetsFound) {
					return SnippetsFetchedMsg([]models.CodeSnippet{})
				} else {
					return ErrMsg{err}
				}
			} else {
				return SnippetsFetchedMsg(snips)
			}
		}
		if m.FilterLang == "" && m.FilterTag == "" {
			snips, err := m.db.GetSnippets(1, 1000)
			if err != nil {
				if errors.Is(err, database.ErrNoSnippetsFound) {
					return SnippetsFetchedMsg([]models.CodeSnippet{})
				} else {
					return ErrMsg{err}
				}
			} else {
				return SnippetsFetchedMsg(snips)
			}
		}
		if m.FilterLang != "" {
			snips, err := m.db.GetSnippetsByLanguage(m.FilterLang)
			if err != nil {
				if errors.Is(err, database.ErrNoSnippetsFound) {
					return SnippetsFetchedMsg([]models.CodeSnippet{})
				} else {
					return ErrMsg{err}
				}
			} else {
				return SnippetsFetchedMsg(snips)
			}
		}
		if m.FilterTag != "" {
			snips, err := m.db.GetSnippetsByTag(m.FilterTag)
			if err != nil {
				if errors.Is(err, database.ErrNoSnippetsFound) {
					return SnippetsFetchedMsg([]models.CodeSnippet{})
				} else {
					return ErrMsg{err}
				}
			} else {
				return SnippetsFetchedMsg(snips)
			}
		}
		return ErrMsg{Err: fmt.Errorf("unable to determine how to get code snippets from db")}
	}
}

func (m Model) generateSnippetsTable(codeSnippets []models.CodeSnippet) tea.Cmd {
	return func() tea.Msg {
		columns := []table.Column{
			{Title: "Name", Width: 25},
			{Title: "Language", Width: 8},
			{Title: "Tags", Width: 20},
			{Title: "Description", Width: 30},
			{Title: "Source", Width: 20},
		}

		var rows []table.Row
		for _, snippet := range codeSnippets {
			rows = append(rows, table.Row{
				snippet.Name, snippet.Language, snippet.Tags, snippet.Description, snippet.Source,
			})
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(m.height-8),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)

		t.SetStyles(s)
		return SnippetsTableMsg(t)
	}
}

func (m Model) openForViewing(codeSnip models.CodeSnippet) tea.Cmd {
	return func() tea.Msg {
		return SingleViewMsg(codeSnip)
	}
}

func (m Model) openForEditing(codeSnip models.CodeSnippet) tea.Cmd {
	return func() tea.Msg {
		return SingleEditMsg(codeSnip)
	}
}
