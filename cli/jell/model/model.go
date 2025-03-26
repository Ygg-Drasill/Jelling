package model

import (
	"os"

	"github.com/Ygg-Drasill/Jelling/cli/jell/model/user"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type JellModel struct {
	user *user.User
}

func InitialModel() JellModel {
	return JellModel{
		user: nil,
	}
}

func (m JellModel) Init() tea.Cmd {
	return nil
}

func (m JellModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			cmd = tea.Quit
			break
		}
	}
	return m, cmd
}

func (m JellModel) View() string {
	termWidth, termHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#553333")).
		PaddingTop(2).
		PaddingLeft(4).
		Height(termHeight).
		Width(termWidth)
	if m.user == nil {
		return style.Render("Not logged in")
	}
	return style.Render("Hello", m.user.Name)
}
