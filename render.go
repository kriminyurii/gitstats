package main

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	commits map[time.Time]int
}

func initialModel() model {
	return model{
		commits: proccessRepos(email),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || s == "q" || s == "esc" {
			return m, tea.Quit
		}

	}
	return m, nil
}

func (m model) View() string {
	var style = lipgloss.NewStyle().
		Width(1).
		Height(1).
		PaddingLeft(1).
		PaddingRight(1).
		SetString("1").
		Background(lipgloss.Color("#FFFFFF"))

	return style.Render()
}

func Render() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
