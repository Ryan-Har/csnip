package editTui

import (
	"github.com/Ryan-Har/csnip/common"
	"github.com/Ryan-Har/csnip/common/models"
	"github.com/Ryan-Har/csnip/database"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	CodeSnippet models.CodeSnippet
	EditMode    bool
	db          database.DatabaseInteractions // used to interact with the database
	TextInputs  []textinput.Model
	//CodeInput    textarea.Model
	focusedInput int // tracks which input is being used. CodeInput is last
	height       int // height of the terminal window
	width        int // width of the terminal window
}

// names editable fields of TextInput, mapped directly to CodeSnippet Fields
const (
	name = iota
	language
	tags
	description
	source
)

const (
	darkGray = lipgloss.Color("#767676")
)

var (
	activeInputStyle = lipgloss.NewStyle().Background(darkGray)
)

// Messages for use in the main model
type ErrMsg struct {
	Err error
}

// manually requests window size
type WindowSizeReqMsg struct{}

func New(db database.DatabaseInteractions, options ...func(*Model)) tea.Model {
	var textInputs []textinput.Model = make([]textinput.Model, 5)

	textInputs[name] = textinput.New()
	textInputs[name].Placeholder = "(Optional) Enter name of the code snippet."
	textInputs[name].CharLimit = 255
	textInputs[name].Prompt = "Name: "
	// validate later
	// textInputs[name].Validate = nameValidator

	textInputs[language] = textinput.New()
	textInputs[language].Placeholder = "Enter the language of the code."
	textInputs[language].CharLimit = 255
	textInputs[language].Prompt = "Language: "
	textInputs[language].ShowSuggestions = true
	textInputs[language].SetSuggestions(common.ListValidLanguages())

	textInputs[tags] = textinput.New()
	textInputs[tags].Placeholder = "(Optional) e.g. Production,Cloudfunctions."
	textInputs[tags].CharLimit = 255
	textInputs[tags].Prompt = "Tags: "

	textInputs[description] = textinput.New()
	textInputs[description].Placeholder = "(Optional) Short description."
	textInputs[description].CharLimit = 255
	textInputs[description].Prompt = "Description: "

	textInputs[source] = textinput.New()
	textInputs[source].Placeholder = "(Optional) Enter a source for the code snippet."
	textInputs[source].CharLimit = 255
	textInputs[source].Prompt = "Source: "

	// codeInput := textarea.New()
	// codeInput.ShowLineNumbers = false

	mod := &Model{
		db:           db,
		CodeSnippet:  models.CodeSnippet{},
		TextInputs:   textInputs,
		focusedInput: 0,
		//CodeInput:   codeInput,
	}
	for _, o := range options {
		o(mod)
	}
	return mod
}

func WithCodeSnippet(snippet models.CodeSnippet) func(*Model) {
	return func(m *Model) {
		m.CodeSnippet = snippet
		m.TextInputs[name].SetValue(snippet.Name)
		m.TextInputs[language].SetValue(snippet.Language)
		m.TextInputs[tags].SetValue(snippet.Tags)
		m.TextInputs[description].SetValue(snippet.Description)
		m.TextInputs[source].SetValue(snippet.Source)
		//m.CodeInput.SetValue(snippet.Code)
	}
}

func InEditMode() func(*Model) {
	return func(m *Model) {
		m.EditMode = true
	}
}

func (m Model) Init() tea.Cmd {
	return m.sendWindowSizeRequest()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ErrMsg:
		panic(msg.Err.Error())
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		for i := range m.TextInputs {
			m.TextInputs[i].Width = m.width - len(m.TextInputs[i].Prompt)
		}
	case tea.KeyMsg:
		cmds = append(cmds, m.sendWindowSizeRequest())
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
	}

	//update text inputs
	for i := range m.TextInputs {
		newModel, newCmd := m.TextInputs[i].Update(msg)
		m.TextInputs[i] = newModel
		cmds = append(cmds, newCmd)
	}
	// newCodeInputModel, newCmd := m.CodeInput.Update(msg)
	// m.CodeInput = newCodeInputModel
	//cmds = append(cmds, newCmd)

	// ensure only a single item if focused, if it's in EditMode
	for i := range m.TextInputs {
		m.TextInputs[i].Blur()
		m.TextInputs[i].TextStyle = lipgloss.NewStyle()
	}
	if m.EditMode {
		m.TextInputs[m.focusedInput].Focus()
		m.TextInputs[m.focusedInput].TextStyle = activeInputStyle
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		m.TextInputs[name].View(),
		m.TextInputs[language].View(),
		m.TextInputs[tags].View(),
		m.TextInputs[description].View(),
		m.TextInputs[source].View(),
	)
	// inputStyle.Width(30).Render("Code"),
	// m.CodeInput.View(),

}

// nextInput focuses the next input field
func (m *Model) nextInput() {
	m.focusedInput = (m.focusedInput + 1) % (len(m.TextInputs) + 1) // + 1
}

// prevInput focuses the previous input field
func (m *Model) prevInput() {
	m.focusedInput--
	// Wrap around
	if m.focusedInput < 0 {
		m.focusedInput = len(m.TextInputs)
	}
}

func (m Model) sendWindowSizeRequest() tea.Cmd {
	return func() tea.Msg {
		return WindowSizeReqMsg{}
	}
}

// type KeyMap struct {
// 	Next key.Binding
// 	Prev key.Binding
// 	Quit key.Binding

// 	Edit key.Binding
// 	Help key.Binding
// }

// var keyMap = KeyMap{
// 	Quit: key.NewBinding(
// 		key.WithKeys("ctrl+c"),
// 		key.WithHelp("ctrl+c", "quit"),
// 	),
// 	Edit: key.NewBinding(
// 		key.WithKeys("e"),
// 		key.WithHelp("e", "edit"),
// 	),
// 	Help: key.NewBinding(
// 		key.WithKeys("?"),
// 		key.WithHelp("?", "toggle help"),
// 	),
// 	Prev: key.NewBinding(
// 		key.WithKeys("shift+tab"),
// 		key.WithHelp("shift+tab", "prev"),
// 	),
// 	Next: key.NewBinding(
// 		key.WithKeys("tab"),
// 		key.WithHelp("tab", "next"),
// 	),
// }

// func (k KeyMap) ShortHelp() []key.Binding {
// 	return []key.Binding{k.Help, k.Quit}
// }

// func (k KeyMap) FullHelp() [][]key.Binding {
// 	return [][]key.Binding{
// 		{k.Up, k.Down, k.View, k.Edit, k.Add},
// 		{k.Help, k.Quit},
// 	}
// }
