package account

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *JellAccountModel) updateInputs(do func(model *textinput.Model)) []tea.Cmd {
	commands := make([]tea.Cmd, len(m.textInputs))
	for i := range m.textInputs {
		if i == m.inputIndex {
			commands[i] = m.textInputs[i].Focus()
			m.textInputs[i].PromptStyle = focusedStyle
			m.textInputs[i].TextStyle = focusedStyle
			if do != nil {
				do(&m.textInputs[i])
			}
			continue
		}
		m.textInputs[i].Blur()
		m.textInputs[i].PromptStyle = noStyle
		m.textInputs[i].TextStyle = noStyle
		if do != nil {
			do(&m.textInputs[i])
		}
	}
	return commands
}

func (m *JellAccountModel) submit() tea.Cmd {
	if m.inputIndex != len(m.textInputs)-1 {
		m.inputIndex++
		return nil
	}

	username := m.textInputs[0].Value()
	password := m.textInputs[1].Value()
	switch m.screenMode {

	case ModeRegister:
		if password != m.textInputs[2].Value() {
			m.resetPasswordInputs()
			return nil
		}
		m.loading = true
		return register(username, password)

	case ModeLogin:
		m.loading = true
		return authenticate(username, password)
	}

	return nil
}

func (m *JellAccountModel) resetPasswordInputs() {
	m.textInputs[1].SetValue("")
	m.textInputs[2].SetValue("")
	m.inputIndex = 1
	m.errorMessage = "Passwords don't match"
}
