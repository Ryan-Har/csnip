package editTui

import (
	"bytes"
	"fmt"

	"github.com/Ryan-Har/csnip/common"
	"github.com/Ryan-Har/csnip/common/models"
	"github.com/Ryan-Har/csnip/database"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type Model struct {
	CodeSnippet  models.CodeSnippet
	EditMode     bool
	db           database.DatabaseInteractions // used to interact with the database
	TextInputs   []textinput.Model
	CodeInput    textarea.Model
	focusedInput int    // tracks which input is being used. CodeInput is last
	keys         KeyMap // KeyMap holding available keys
	help         help.Model
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

// Messages for use in the main model
type ErrMsg struct {
	Err error
}

// manually requests window size
type WindowSizeReqMsg struct{}
type ReturnToViewMsg struct{}

func New(db database.DatabaseInteractions, options ...func(*Model)) tea.Model {
	var textInputs []textinput.Model = make([]textinput.Model, 5)

	textInputs[name] = textinput.New()
	textInputs[name].Placeholder = "(Optional) Enter name of the code snippet."
	textInputs[name].CharLimit = 255
	textInputs[name].Prompt = "Name: "
	textInputs[name].TextStyle = lipgloss.NewStyle()
	textInputs[name].TextStyle.Height(1)
	// validate later
	// textInputs[name].Validate = nameValidator

	textInputs[language] = textinput.New()
	textInputs[language].Placeholder = "Enter the language of the code."
	textInputs[language].CharLimit = 255
	textInputs[language].Prompt = "Language: "
	textInputs[language].ShowSuggestions = true
	textInputs[language].SetSuggestions(common.ListValidLanguages())
	textInputs[language].TextStyle = lipgloss.NewStyle()

	textInputs[tags] = textinput.New()
	textInputs[tags].Placeholder = "(Optional) e.g. Production,Cloudfunctions."
	textInputs[tags].CharLimit = 255
	textInputs[tags].Prompt = "Tags: "
	textInputs[tags].TextStyle = lipgloss.NewStyle()

	textInputs[description] = textinput.New()
	textInputs[description].Placeholder = "(Optional) Short description."
	textInputs[description].CharLimit = 255
	textInputs[description].Prompt = "Description: "

	textInputs[source] = textinput.New()
	textInputs[source].Placeholder = "(Optional) Enter a source for the code snippet."
	textInputs[source].CharLimit = 255
	textInputs[source].Prompt = "Source: "

	codeInput := textarea.New()
	codeInput.ShowLineNumbers = false
	codeInput.Prompt = ""
	codeInput.CharLimit = 0 // no limit on code input

	mod := &Model{
		db:           db,
		CodeSnippet:  models.CodeSnippet{},
		TextInputs:   textInputs,
		keys:         keyMap,
		help:         help.New(),
		focusedInput: 0,
		CodeInput:    codeInput,
	}
	for _, o := range options {
		o(mod)
	}
	if mod.EditMode {
		mod.keys.Edit.SetEnabled(false)
		mod.keys.Save.SetEnabled(true)
	} else {
		mod.keys.Edit.SetEnabled(true)
		mod.keys.Save.SetEnabled(false)
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
		m.CodeInput.SetValue(snippet.Code)
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
		m.CodeInput.SetWidth(m.width)
		// 9 is the total number of lines used by the help dialogue and textInput fields. This allows the codeInput to use as much space vertically as it can.
		m.CodeInput.SetHeight(m.height - 9)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Edit):
			m.EditMode = true
			m.keys.Edit.SetEnabled(false)
			m.keys.Save.SetEnabled(true)
		case key.Matches(msg, m.keys.Save):
			return m, m.saveToDb()
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Next):
			m.nextInput()
		case key.Matches(msg, m.keys.Prev):
			m.prevInput()
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		default:
			//update text inputs if it isn't one of the above keys
			for i := range m.TextInputs {
				newModel, newCmd := m.TextInputs[i].Update(msg)
				m.TextInputs[i] = newModel
				cmds = append(cmds, newCmd)
			}
			var cmd tea.Cmd
			m.CodeInput, cmd = m.CodeInput.Update(msg)
			cmds = append(cmds, cmd)
		}
		// explicit convert of tabs to spaces
		if msg.Type == tea.KeyTab && m.CodeInput.Focused() && m.EditMode {
			m.CodeInput.InsertString("    ")
		}

	}

	// ensure only a single item if focused, if it's in EditMode
	// blur all inputs each time
	for i := range m.TextInputs {
		m.TextInputs[i].Blur()
		m.TextInputs[i].TextStyle = lipgloss.NewStyle()
	}
	m.CodeInput.Blur()

	// unblur and focus the currently edited item
	if m.EditMode {
		if m.focusedInput == len(m.TextInputs) {
			m.CodeInput.Focus()
		} else {
			m.TextInputs[m.focusedInput].Focus()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	helpview := m.help.View(m.keys)

	return lipgloss.JoinVertical(lipgloss.Top,
		m.TextInputs[name].View(),
		m.TextInputs[language].View(),
		m.TextInputs[tags].View(),
		m.TextInputs[description].View(),
		m.TextInputs[source].View(),
		"Code:",
		//TODO: handle this correctly so that syntax highlighting works with the method
		m.CodeInput.View(),
		//highlightCode(m.CodeInput.Value(), m.TextInputs[language].Value()),
	) + "\n" + helpview
}

// nextInput focuses the next input field
func (m *Model) nextInput() {
	m.focusedInput = (m.focusedInput + 1) % (len(m.TextInputs) + 1)
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

func (m Model) saveToDb() tea.Cmd {
	return func() tea.Msg {
		m.CodeSnippet.Code = m.CodeInput.Value()
		m.CodeSnippet.Name = m.TextInputs[name].Value()
		m.CodeSnippet.Language = m.TextInputs[language].Value()
		m.CodeSnippet.Tags = m.TextInputs[tags].Value()
		m.CodeSnippet.Description = m.TextInputs[description].Value()
		m.CodeSnippet.Source = m.TextInputs[source].Value()

		if m.CodeSnippet.Uuid == uuid.Nil {
			if err := m.db.AddNewSnippet(m.CodeSnippet); err != nil {
				return ErrMsg{Err: fmt.Errorf("unable to save to database %v", err)}
			}
		} else {
			if _, err := m.db.UpdateSnippet(m.CodeSnippet.Uuid, m.CodeSnippet); err != nil {
				return ErrMsg{Err: fmt.Errorf("unable to save to database %v", err)}
			}
		}
		return ReturnToViewMsg{}
	}
}

type KeyMap struct {
	Next key.Binding
	Prev key.Binding
	Quit key.Binding

	Edit key.Binding
	Help key.Binding
	Save key.Binding
}

var keyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Help: key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "toggle help"),
	),
	Prev: key.NewBinding(
		key.WithKeys("ctrl+b"),
		key.WithHelp("ctrl+b", "prev"),
	),
	Next: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "next"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Next, k.Prev, k.Edit, k.Save},
		{k.Help, k.Quit},
	}
}

func highlightCode(code string, lang string) string {
	themeString := "monokai"
	if code != "" && lang != "" {
		lexer := lexers.Get(lang)
		if lexer == nil {
			lexer = lexers.Fallback
		}
		style := styles.Get(themeString)
		formatter := formatters.Get("terminal256")
		if formatter == nil {
			formatter = formatters.Fallback
		}
		iterator, _ := lexer.Tokenise(nil, code)

		buf := new(bytes.Buffer)
		_ = formatter.Format(buf, style, iterator)
		fmt.Println(buf.String())
		return buf.String()
	}
	return code
}
