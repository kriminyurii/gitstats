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

type Styles struct {
	Cell           lipgloss.Style
	MonthsRow      lipgloss.Style
	Month          lipgloss.Style
	Row            lipgloss.Style
	RowStartIndent int
	MonthsIndent   int
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.MonthsRow = lipgloss.NewStyle().MarginTop(1).MarginBottom(1)
	s.Month = lipgloss.NewStyle().Bold(true)
	s.Cell = lipgloss.NewStyle().
		Width(1).
		Height(1).
		PaddingRight(1).
		PaddingLeft(1).
		Background(lipgloss.Color("#FFFFFF"))
	s.Row = lipgloss.NewStyle().MarginBottom(1)
	s.RowStartIndent = 4
	s.MonthsIndent = 13

	return s
}

func RenderMonthsRow(s string) string {
	styles := DefaultStyles()
	months := GetLastHalfYearInMonths()
	monthsRow := ""
	for i, mon := range months {
		if i == 0 {
			monthsRow += styles.Month.Render(mon)
		} else {
			monthsRow += styles.Month.MarginLeft(styles.MonthsIndent).Render(mon)
		}
	}
	s += styles.MonthsRow.MarginLeft(styles.RowStartIndent).Render(monthsRow)
	s += "\n"
	return s
}

func RenderRow(s string) string {
	styles := DefaultStyles()
	daysOfWweek := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	weeksInMonth := 4
	monthsCount := len(GetLastHalfYearInMonths())

	for i, day := range daysOfWweek {
		row := ""
		if i%2 != 0 {
			row = lipgloss.JoinHorizontal(lipgloss.Center, day, styles.Cell.MarginLeft(1).Render("0"))
			for i := 1; i < weeksInMonth*monthsCount; i++ {
				row += lipgloss.JoinHorizontal(lipgloss.Center, styles.Cell.MarginLeft(1).Render("0"))
			}
		} else if i%2 == 0 {
			row = styles.Cell.MarginLeft(styles.RowStartIndent).Render("0")
			for i := 1; i < weeksInMonth*monthsCount; i++ {
				row += lipgloss.JoinHorizontal(lipgloss.Center, styles.Cell.MarginLeft(1).Render("0"))
			}
		}
		s += styles.Row.Render(row)
		s += "\n"
	}
	return s
}

func (m model) View() string {
	s := ""
	s += RenderMonthsRow(s)
	s += RenderRow(s)
	return s
}

func Render() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
