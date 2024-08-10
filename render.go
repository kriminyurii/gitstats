package main

import (
	"log"
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

var commitsLevels = []struct {
	threshold int
	color     string
}{
	{0, "#6F7573"}, // No contributions
	{1, "#C6E48B"},
	{5, "#7BC96F"},
	{10, "#239A3B"},
	{20, "#196127"},
}

func getColor(commits int) string {
	for _, level := range commitsLevels {
		if commits <= level.threshold {
			return level.color
		}
	}
	return commitsLevels[len(commitsLevels)-1].color
}

type Styles struct {
	Cell           lipgloss.Style
	MonthsRow      lipgloss.Style
	Month          lipgloss.Style
	Row            lipgloss.Style
	Info           lipgloss.Style
	RowStartIndent int
	MonthsIndent   int
	InfoIndent     int
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.MonthsRow = lipgloss.NewStyle().MarginTop(1).MarginBottom(1)
	s.Month = lipgloss.NewStyle().Bold(true)
	s.Cell = lipgloss.NewStyle().
		Width(3).
		Height(1).
		PaddingRight(1).
		PaddingLeft(1).
		Background(lipgloss.Color("#FFFFFF"))
	s.Row = lipgloss.NewStyle().MarginBottom(1)
	s.RowStartIndent = 4
	s.MonthsIndent = 13
	monthsCount := len(GetLastHalfYearInMonths())
	infoWordsIndent := 8
	infoIndent := s.RowStartIndent + s.MonthsIndent*monthsCount - 1 - infoWordsIndent
	s.InfoIndent = infoIndent
	s.Info = lipgloss.NewStyle().MarginLeft(s.InfoIndent).Bold(true)

	return s
}

func getCellValue(weekdayIndex, rowIndex int, model model) int {
	var commitsCount int
	for _, offset := range model.commits {
		if weekdayIndex == offset.WeekDay && rowIndex == offset.Row {
			commitsCount = offset.Commits
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

func renderCell(styles *Styles, cellValue int) lipgloss.Style {
	cell := styles.Cell
	color := getColor(cellValue)
	cell = cell.Background(lipgloss.Color(color))
	return cell
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
			row = lipgloss.JoinHorizontal(lipgloss.Center, day, renderCell(styles, cellValue).MarginLeft(1).Render())
			for j := 1; j < weeksInMonth*monthsCount; j++ {
				cellValue = getCellValue(i, j, model)
				row += lipgloss.JoinHorizontal(lipgloss.Center, renderCell(styles, cellValue).MarginLeft(1).Render())
			}
		} else if i%2 == 0 {
			row = renderCell(styles, cellValue).MarginLeft(styles.RowStartIndent).Render()
			for j := 1; j < weeksInMonth*monthsCount; j++ {
				cellValue = getCellValue(i, j, model)
				row += lipgloss.JoinHorizontal(lipgloss.Center, renderCell(styles, cellValue).MarginLeft(1).Render())
			}
		}
		renderString += styles.Row.Render(row)

		renderString += "\n"

	}
	return renderString
}

func RenderCommitsLevelInfo() string {
	styles := DefaultStyles()
	infoRow := styles.Info
	info := ""
	info += lipgloss.JoinHorizontal(lipgloss.Center, "Less")
	for _, level := range commitsLevels {
		info += lipgloss.JoinHorizontal(lipgloss.Center, renderCell(styles, level.threshold).MarginLeft(1).Render())
	}
	info += lipgloss.JoinHorizontal(lipgloss.Center, lipgloss.NewStyle().MarginLeft(1).SetString("More").Render())
	return infoRow.Render(info)
}

func (m model) View() string {
	s := ""
	s += RenderMonthsRow()
	s += "\n"
	s += RenderGrid(m)
	s += RenderCommitsLevelInfo()
	return s
}

func Render() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
