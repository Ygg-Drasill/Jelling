package ui

import (
	svc "github.com/Ygg-Drasill/Jelling/service/client"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"os"
)

type MenuItem struct {
	title string
	desc  string
}

func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.desc }
func (i MenuItem) FilterValue() string { return i.title }

type JellState struct {
	Menu  list.Model
	State textinput.Model
}

func NewJellState() JellState {
	items := []list.Item{
		MenuItem{title: "Search", desc: "Find articles"},
		MenuItem{title: "Github", desc: "Fetch file from Github"},
		MenuItem{title: "Exit", desc: "Quit the application"},
	}

	delegate := list.NewDefaultDelegate()
	termWidth, termHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	menu := list.New(items, delegate, termWidth, termHeight)
	menu.Title = "Main Menu"
	menu.SetFilteringEnabled(false)

	return JellState{Menu: menu}
}

func SearchState() JellState {
	searchInput := textinput.New()
	searchInput.Placeholder = "Search"
	searchInput.Focus()
	searchInput.CharLimit = 256
	searchInput.Width = 30

	searchInput.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	searchInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4287f5"))

	return JellState{State: searchInput}
}

func GithubState() JellState {
	input := textinput.New()
	input.Prompt = string(owner)
	input.Placeholder = "owner/repo/filepath"
	input.Focus()
	input.CharLimit = 256
	input.Width = 30
	input.ShowSuggestions = true
	suggestions := []string{"Ygg-Drasill", "Jelling", "README.md"}
	input.SetSuggestions(suggestions)

	input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))
	input.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4287f5"))

	return JellState{State: input}
}

func (js JellState) HandleEnter(component any) tea.Cmd {
	switch component := component.(type) {
	case textinput.Model:
		component.Blur()
	}

	return nil
}

var v []string

type githubPrompt string

const (
	owner    githubPrompt = "Enter the owner: "
	repo     githubPrompt = "Enter the repo: "
	filepath githubPrompt = "Enter the filepath: "
)

func (js JellState) HandleGitHub(t textinput.Model) textinput.Model {
	if t.Value() != "" {
		switch t.Prompt {
		case string(owner):
			v = append(v, t.Value())
			t.Prompt = string(repo)
			t.SetValue("")
			return t
		case string(repo):
			v = append(v, t.Value())
			t.Prompt = string(filepath)
			t.SetValue("")
			return t
		case string(filepath):
			v = append(v, t.Value())
			d := svc.FetchFromGitHub(v)
			v = []string{}
			t.Prompt = string(d)
			t.SetValue("")
			t.Placeholder = ""
			return t
		}
	}
	if t.Placeholder == "" {
		t.Prompt = string(owner)
		t.Placeholder = "owner/repo/filepath"
		return t
	}
	return t
}
