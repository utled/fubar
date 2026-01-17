package tui

import (
	"fmt"
	"fubar/data"
	"fubar/utils"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type statsCollectedMsg struct {
	monthStats data.MonthStats
	fullStats  data.FullStats
}

func (model *Model) collectStats() tea.Msg {
	return func() tea.Msg {
		return statsCollectedMsg{}
	}
}

type statsYearMsg struct {
	minYear int
	maxYear int
}

func (model *Model) collectYearRange() tea.Cmd {
	return func() tea.Msg {
		minYear, maxYear, err := data.GetYearRange()
		if err != nil {
			return errMsg(err)
		}
		return statsYearMsg{
			minYear: minYear,
			maxYear: maxYear,
		}
	}
}

type statsGraphMsg struct {
	graphString string
}

func (model *Model) generateStatsGraph() tea.Cmd {
	return func() tea.Msg {
		selectedYear, err := strconv.Atoi(model.statsDetails.yearSelection.Value())
		if err != nil {
			return errMsg(err)
		}
		monthlySummary, err := data.GetMonthlySummary(selectedYear)
		if err != nil {
			return err
		}

		maxDays := 25
		// If table height and legend should only reach max worked/week days
		/*var maxDays int
		for _, month := range monthlySummary {
			if month.WorkedDays > maxDays {
				maxDays = month.WorkedDays
			}
			if month.TotalWeekDays > maxDays {
				maxDays = month.TotalWeekDays
			}
		}*/

		var b strings.Builder

		// Table Header
		fmt.Fprintf(&b, "%-5s(Worked Days ██, Weekdays ░)\n", "")

		// Column Titles (Month names)
		months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
		fmt.Fprintf(&b, "%-5s", "")
		for _, month := range months {
			fmt.Fprintf(&b, "%-4s", month)
		}
		b.WriteString("\n")

		// Column Subtitles (count of worked days that month)
		fmt.Fprintf(&b, "%-5s", "")
		for idx := range months {
			val := 0
			if len(monthlySummary) >= idx+1 {
				val = monthlySummary[idx].WorkedDays
			}
			fmt.Fprintf(&b, "%-4d", val)
		}
		b.WriteString("\n")

		// Top border
		fmt.Fprintf(&b, "%-4s", "")
		b.WriteString(strings.Repeat("_", 48) + "\n")

		// Table rows
		for i := maxDays - 1; i >= 0; i-- {
			// Left border
			fmt.Fprintf(&b, "%-3d|", i+1)
			for j := 0; j < 12; j++ {
				if len(monthlySummary) >= j+1 && monthlySummary[j].WorkedDays >= i+1 {
					fmt.Fprintf(&b, "%3s", "██")
				} else {
					fmt.Fprintf(&b, "%3s", "  ")
				}
				if len(monthlySummary) >= j+1 && monthlySummary[j].TotalWeekDays >= i+1 {
					b.WriteString("░")
				} else {
					b.WriteString(" ")
				}
			}
			// Right border
			b.WriteString("|\n")
		}

		drawnGraph := b.String()
		return statsGraphMsg{graphString: drawnGraph}
	}
}

type statsTableDataMsg struct {
	tableData []*data.MonthStats
}

func (model *Model) fetchStatsTableData() tea.Cmd {
	return func() tea.Msg {
		selectedYear, err := strconv.Atoi(model.statsDetails.yearSelection.Value())
		monthlySummary, err := data.GetMonthlySummary(selectedYear)
		if err != nil {
			return err
		}
		return statsTableDataMsg{tableData: monthlySummary}
	}
}

type statsAllSumDataMsg struct {
	fieldData *data.FullStats
}

func (model *Model) fetchAllSumData() tea.Cmd {
	return func() tea.Msg {
		minDate := "2000-01-01"
		maxDate := model.dateState.MaxDate
		fullStatistics, err := data.GetFullStatistics(minDate, maxDate)
		if err != nil {
			return err
		}
		return statsAllSumDataMsg{fieldData: fullStatistics}
	}
}

type statsMonthSumDataMsg struct {
	fieldData *data.FullStats
}

func (model *Model) fetchMonthSumData(monthName string) tea.Cmd {

	return func() tea.Msg {
		var startDate, endDate string
		yearInt, err := strconv.Atoi(model.statsDetails.yearSelection.Value())
		if err != nil {
			return errMsg(err)
		}

		switch monthName {
		case "":
		case "Total":
			startDate = time.Date(yearInt, 1, 1, 0, 0, 0, 0, time.UTC).Format(utils.DateLayout)
			firstOfNextJanuary := time.Date(yearInt+1, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate = firstOfNextJanuary.Add(-time.Nanosecond).Format(utils.DateLayout)
		default:
			month, err := time.Parse("January", monthName)
			if err != nil {
				return errMsg(err)
			}
			monthInt := int(month.Month())
			startDate = time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC).Format(utils.DateLayout)
			firstOfNextMonth := time.Date(yearInt, time.Month(monthInt)+1, 1, 0, 0, 0, 0, time.UTC)
			endDate = firstOfNextMonth.Add(-time.Nanosecond).Format(utils.DateLayout)
		}

		fullStatistics, err := data.GetFullStatistics(startDate, endDate)
		if err != nil {
			return err
		}
		return statsMonthSumDataMsg{fieldData: fullStatistics}
	}
}
