package account

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type JellAccountModel struct {
	screenMode   ScreenMode
	textInputs   []textinput.Model
	inputIndex   int
	loading      bool
	errorMessage string
}

type ScreenMode int

const (
	ModeLogin    ScreenMode = iota
	ModeRegister ScreenMode = iota
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("179"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("161"))
)

func InitAccountModel(mode ScreenMode) JellAccountModel {
	username := textinput.New()
	password := textinput.New()
	username.Placeholder = "Username"
	password.Placeholder = "Password"

	username.Focus()
	username.PromptStyle = focusedStyle
	username.TextStyle = focusedStyle
	username.Cursor.Style = focusedStyle

	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = 'ᛜ'
	password.Cursor.Style = focusedStyle

	textInputs := make([]textinput.Model, 2)
	textInputs[0] = username
	textInputs[1] = password

	if mode == ModeRegister {
		passwordConfirm := textinput.New()
		passwordConfirm.Placeholder = "Confirm Password"
		passwordConfirm.EchoMode = textinput.EchoPassword
		passwordConfirm.EchoCharacter = 'ᛜ'
		textInputs = append(textInputs, passwordConfirm)
	}

	for i := range textInputs { //TODO: fix
		textInputs[i].PlaceholderStyle.Foreground(lipgloss.Color("179"))
	}

	return JellAccountModel{
		screenMode: mode,
		textInputs: textInputs,
		inputIndex: 0,
	}
}

func (m JellAccountModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m JellAccountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updateCommands := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.inputIndex == len(m.textInputs)-1 {
				if m.textInputs[1].Value() != m.textInputs[2].Value() {
					m.textInputs[1].SetValue("")
					m.textInputs[2].SetValue("")
					m.inputIndex = 1
					m.errorMessage = "Passwords don't match"
				} else {
					m.loading = true
					return m, tea.Batch(updateCommands...)
				}
			} else {
				m.inputIndex++
			}

			commands := make([]tea.Cmd, len(m.textInputs))
			for i := range m.textInputs {
				if i == m.inputIndex {
					commands[i] = m.textInputs[i].Focus()
					m.textInputs[i].PromptStyle = focusedStyle
					m.textInputs[i].TextStyle = focusedStyle
					continue
				}
				m.textInputs[i].Blur()
				m.textInputs[i].PromptStyle = noStyle
				m.textInputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(commands...)
		}
	}

	if m.inputIndex >= len(m.textInputs) {
		m.inputIndex = len(m.textInputs)
	} else if m.inputIndex < 0 {
		m.inputIndex = 0
	}

	for i := range m.textInputs {
		var inputCommand tea.Cmd
		m.textInputs[i], inputCommand = m.textInputs[i].Update(msg)
		updateCommands = append(updateCommands, inputCommand)
	}

	return m, tea.Batch(updateCommands...)
}

func (m JellAccountModel) View() string {
	if m.loading {
		return focusedStyle.Render("loading...")
	}

	var builder strings.Builder

	if m.errorMessage != "" {
		builder.WriteString(errorStyle.Render(m.errorMessage))
		builder.WriteRune('\n')
	}

	for _, textInput := range m.textInputs {
		builder.WriteString(textInput.View())
		builder.WriteRune('\n')
	}
	return builder.String()
}
