package editTui

import (
	"github.com/Ryan-Har/csnip/common/models"
	"github.com/Ryan-Har/csnip/database"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	CodeSnippet models.CodeSnippet
	EditMode    bool
	db          database.DatabaseInteractions // used to interact with the database
	TextEditor  string
}

func New(db database.DatabaseInteractions, options ...func(*Model)) tea.Model {
	mod := &Model{
		db: db,
	}
	for _, o := range options {
		o(mod)
	}
	return mod
}

func WithCodeSnippet(snippet models.CodeSnippet) func(*Model) {
	return func(m *Model) {
		m.CodeSnippet = snippet
	}
}

func InEditMode() func(*Model) {
	return func(m *Model) {
		m.EditMode = true
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return m, nil
}

func (m Model) View() string {
	return "Single View" + "\n"
}
