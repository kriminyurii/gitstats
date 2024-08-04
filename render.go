package main

import (
	"log"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	commits map[time.Time]Offset
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

func getCellValue(weekdayIndex, rowIndex int, model model) string {
	commitsCount := "0"
	for _, offset := range model.commits {
		if weekdayIndex == offset.WeekDay && rowIndex == offset.Row {
			commitsCount = strconv.Itoa(offset.Commits)
		}
	}
	return commitsCount
}

func RenderMonthsRow() string {
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
	return styles.MonthsRow.MarginLeft(styles.RowStartIndent).Render(monthsRow)
}

func RenderGrid(model model) string {
	var renderString string = ""
	styles := DefaultStyles()
	daysOfWweek := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	weeksInMonth := 4
	monthsCount := len(GetLastHalfYearInMonths())

	for i, day := range daysOfWweek {
		row := ""
		cellValue := getCellValue(i, 0, model)
		if i%2 != 0 {
			row = lipgloss.JoinHorizontal(lipgloss.Center, day, styles.Cell.MarginLeft(1).Render(cellValue))
			for j := 1; j < weeksInMonth*monthsCount; j++ {
				cellValue = getCellValue(i, j, model)
				row += lipgloss.JoinHorizontal(lipgloss.Center, styles.Cell.MarginLeft(1).Render(cellValue))
			}
		} else if i%2 == 0 {
			row = styles.Cell.MarginLeft(styles.RowStartIndent).Render(cellValue)
			for j := 1; j < weeksInMonth*monthsCount; j++ {
				cellValue = getCellValue(i, j, model)
				row += lipgloss.JoinHorizontal(lipgloss.Center, styles.Cell.MarginLeft(1).Render(cellValue))
			}
		}
		renderString += styles.Row.Render(row)
		renderString += "\n"
	}
	return renderString
}

func (m model) View() string {
	s := ""
	s += RenderMonthsRow()
	s += "\n"
	s += RenderGrid(m)
	s += "\n"
	return s
}

func Render() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
