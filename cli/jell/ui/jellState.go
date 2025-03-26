package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type JellState struct {
	SearchInput textinput.Model
}

func NewJellState() JellState {
	searchInput := textinput.New()
	searchInput.Placeholder = "Search"
	searchInput.Focus()
	searchInput.CharLimit = 256
	searchInput.Width = 30

	searchInput.Prompt = "> "

	return JellState{
		SearchInput: searchInput,
	}
}

func (js JellState) HandleEnter(component any) tea.Cmd {
	switch component := component.(type) {
	case textinput.Model:
		component.Blur()
	}

	return nil
}
