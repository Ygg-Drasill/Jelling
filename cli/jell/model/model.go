package model

import (
	"github.com/Ygg-Drasill/Jelling/cli/jell/ui"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type JellModel struct {
	state State
	menu  list.Model
	input textinput.Model
}

func InitialModel() JellModel {
	s := ui.NewJellState()
	return JellModel{
		state: menu,
		menu:  s.Menu,
	}
}

func (m JellModel) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Jelling"),
		textinput.Blink,
	)
}

func (m JellModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case menu:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter":
				if selectedItem, ok := m.menu.SelectedItem().(ui.MenuItem); ok {
					switch strings.ToLower(selectedItem.Title()) {
					case string(search):
						m.state = search
						m.input = ui.SearchState().State
						return m, textinput.Blink
					case string(github):
						m.state = github
						m.input = ui.GithubState().State
						return m, textinput.Blink
					case "Exit":
						return m, tea.Quit
					}
				}
			}
			m.menu, cmd = m.menu.Update(msg)
		case search:
			switch msg.String() {
			case "esc", "ctrl+c":
				m.state = menu
				return m, nil
			case "enter":
				ui.SearchState().HandleEnter(m.input)
			}
			m.input, cmd = m.input.Update(msg)
		case github:
			switch msg.String() {
			case "esc", "ctrl+c":
				m.state = menu
				return m, nil
			case "enter":
				m.input = ui.GithubState().HandleGitHub(m.input)
			}
			m.input, cmd = m.input.Update(msg)
		}
	}
	return m, cmd
}

func (m JellModel) View() string {
	switch m.state {
	case menu:
		return m.menu.View()
	case search:
		return m.input.View()
	case github:
		return m.input.View()
	}
	return "Hello"
}
