package model

import (
	"os"

	"github.com/Ygg-Drasill/Jelling/cli/jell/ui"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type JellModel struct {
	user            *User
	articleResults  []Article
	ui              ui.JellState
	activeComponent any
}

func InitialModel() JellModel {
	state := ui.NewJellState()
	return JellModel{
		user:            nil,
		activeComponent: state.SearchInput,
		ui:              state,
	}
}

func (m JellModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m JellModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			m.activeComponent = nil
			cmds = append(cmds, m.ui.HandleEnter(m.activeComponent))
		}
		if m.activeComponent != nil {
			break
		}
		switch msg.String() {
		case "q", "ctrl+c":
			cmds = append(cmds, tea.Quit)
			break
		}
	}

	var searchCmd tea.Cmd
	m.ui.SearchInput, searchCmd = m.ui.SearchInput.Update(msg)
	cmds = append(cmds, searchCmd)
	return m, tea.Batch(cmds...)
}

func (m JellModel) View() string {
	termWidth, termHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	style := lipgloss.NewStyle().
		Height(termHeight).
		Width(termWidth).
		Align(lipgloss.Center)
	if m.user == nil {
		return style.Render(m.ui.SearchInput.View())
	}
	return style.Render("Hello", m.user.Name)
}
