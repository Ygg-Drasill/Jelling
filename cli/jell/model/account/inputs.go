package account

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m JellAccountModel) updateInputs(do func(model *textinput.Model)) []tea.Cmd {
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
