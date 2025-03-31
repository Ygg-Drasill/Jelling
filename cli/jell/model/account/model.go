package account

import (
	"fmt"
	"github.com/Ygg-Drasill/Jelling/cli/jell/ui"
	"github.com/charmbracelet/bubbles/spinner"
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
	spinner      spinner.Model
	errorMessage string
	done         bool
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

	spinnerModel := spinner.New()
	spinnerModel.Spinner = spinner.Dot
	spinnerModel.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(ui.Theme.Primary))

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
		spinner:    spinnerModel,
	}
}

func (m JellAccountModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m JellAccountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updateCommands := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			updateCommands = append(updateCommands, m.submit())
			updateCommands = append(updateCommands, m.updateInputs(nil)...)
		}

	case FetchCompleteMsg:
		m.loading = false
		if msg.err != nil {
			m.errorMessage = msg.err.Error()
			m.inputIndex = 0
			m.updateInputs(func(input *textinput.Model) {
				input.SetValue("")
			})
		} else {
			m.done = true
			return m, tea.Quit
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

	var spinnerCommand tea.Cmd
	m.spinner, spinnerCommand = m.spinner.Update(msg)
	updateCommands = append(updateCommands, spinnerCommand)

	return m, tea.Batch(updateCommands...)
}

func (m JellAccountModel) View() string {
	if m.done {
		return ""
	}
	if m.loading {
		return fmt.Sprintf("%s Loading...", m.spinner.View())
	}

	var builder strings.Builder

	if len(m.errorMessage) > 0 {
		builder.WriteString(errorStyle.Render(m.errorMessage))
		builder.WriteRune('\n')
	}

	for _, textInput := range m.textInputs {
		builder.WriteString(textInput.View())
		builder.WriteRune('\n')
	}
	return builder.String()
}
