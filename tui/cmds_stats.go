package tui

import (
	"fmt"
	"fubar/data"
	"fubar/helpers"
	"fubar/utils"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
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
		// If graph height and legend should only reach max worked/week days
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

		// Graph rows
		for i := maxDays - 1; i >= 0; i-- {
			// Left border
			fmt.Fprintf(&b, "%-3d|", i+1)
			for j := range 12 {
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
			b.WriteString(" |\n")
		}

		// Footer
		b.WriteString("( Worked Days ██, Weekdays ░)")

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
		if err != nil {
			return err
		}
		monthlySummary, err := data.GetMonthlySummary(selectedYear)
		if err != nil {
			return err
		}
		return statsTableDataMsg{tableData: monthlySummary}
	}
}

type statsYearSumDataMsg struct {
	fieldData *data.FullStats
}

func (model *Model) fetchYearSumData() tea.Cmd {
	return func() tea.Msg {
		minDate := model.statsDetails.yearSelection.Value() + "-01-01"
		maxDate := model.statsDetails.yearSelection.Value() + "-12-31"
		fullStatistics, err := data.GetFullStatistics(minDate, maxDate)
		if err != nil {
			return err
		}
		return statsYearSumDataMsg{fieldData: fullStatistics}
	}
}

type statsMonthAvgDataMsg struct {
	fieldData *data.FullStats
}

func (model *Model) fetchMonthAvgData(monthName string) tea.Cmd {

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
		return statsMonthAvgDataMsg{fieldData: fullStatistics}
	}
}

func generateTableRow(record *data.WorkDateRecord) table.Row {
	workDate := record.WorkDate
	rowDayType := record.DayType.String
	var startTime string
	if record.StartTime.Valid {
		startTime = record.StartTime.String[:5]
	} else {
		startTime = ""
	}
	lunchDuration := strconv.Itoa(int(record.LunchDuration.Int16))
	var endTime string
	if record.EndTime.Valid {
		endTime = record.EndTime.String[:5]
	} else {
		endTime = ""
	}
	additionalTime := strconv.Itoa(int(record.AdditionalTime.Int16))
	var dayTotal string
	if record.DayTotal.Valid {
		dayTotal = record.DayTotal.String[:5]
	} else {
		dayTotal = ""
	}
	var overtime string
	if record.Overtime.Valid {
		overtime = strconv.FormatBool(record.Overtime.Bool)
	} else {
		overtime = ""
	}
	var dayBalance string
	if record.DayBalance.Valid {
		dayBalance = fmt.Sprintf("%6s", helpers.DecimalToTime(record.DayBalance.Float64))
	} else {
		dayBalance = ""
	}
	var totalBalance string
	if record.TotalBalance.Valid {
		totalBalance = fmt.Sprintf("%6s", helpers.DecimalToTime(record.TotalBalance.Float64))
	} else {
		totalBalance = ""
	}
	tableRow := table.Row{
		workDate,
		rowDayType,
		startTime,
		lunchDuration,
		endTime,
		additionalTime,
		dayTotal,
		overtime,
		dayBalance,
		totalBalance,
	}

	return tableRow
}

func generateStatsTableRow(monthlySummary *data.MonthStats, monthlyTotals *data.MonthStats) table.Row {
	monthlyTotals.TotalWeekDays += monthlySummary.TotalWeekDays
	monthlyTotals.WorkedDays += monthlySummary.WorkedDays

	existingTimes := strings.Split(monthlyTotals.WorkedTime, ":")
	var existingHours, existingMinutes int
	var err error
	if len(existingTimes) > 1 {
		existingHours, err = strconv.Atoi(existingTimes[0])
		if err != nil {
			existingHours = 0
		}
		existingMinutes, err = strconv.Atoi(existingTimes[1])
		if err != nil {
			existingMinutes = 0
		}
	}

	currentTimes := strings.Split(monthlySummary.WorkedTime, ":")
	var currentHours, currentMinutes int
	if len(currentTimes) > 1 {
		currentHours, err = strconv.Atoi(currentTimes[0])
		if err != nil {
			currentHours = 0
		}
		currentMinutes, err = strconv.Atoi(currentTimes[1])
		if err != nil {
			currentMinutes = 0
		}
	}

	hoursBefore := existingHours + currentHours
	minutesBefore := existingMinutes + currentMinutes
	totalHours := hoursBefore + (minutesBefore / 60)
	totalMinutes := minutesBefore % 60
	var suffix string
	if totalMinutes == 0 {
		suffix = "0"
	}

	monthlyTotals.WorkedTime = fmt.Sprintf("%d:%d%s", totalHours, totalMinutes, suffix)
	monthlyTotals.VacationDays += monthlySummary.VacationDays
	monthlyTotals.SickDays += monthlySummary.SickDays
	monthlyTotals.WeekendDays += monthlySummary.WeekendDays
	monthlyTotals.OffDays += monthlySummary.OffDays
	monthlyTotals.OverTimeDays += monthlySummary.OverTimeDays
	monthlyTotals.TotalOvertime.Float64 += monthlySummary.TotalOvertime.Float64

	totalOvertime := helpers.DecimalToTime(monthlySummary.TotalOvertime.Float64)
	workedTime := strings.Split(monthlySummary.WorkedTime, ":")
	tableRow := table.Row{
		monthlySummary.Month,
		strconv.Itoa(monthlySummary.TotalWeekDays),
		strconv.Itoa(monthlySummary.WorkedDays),
		workedTime[0] + ":" + workedTime[1],
		strconv.Itoa(monthlySummary.VacationDays),
		strconv.Itoa(monthlySummary.SickDays),
		strconv.Itoa(monthlySummary.WeekendDays),
		strconv.Itoa(monthlySummary.OffDays),
		strconv.Itoa(monthlySummary.OverTimeDays),
		totalOvertime,
	}

	return tableRow
}

func generateBlankRow() table.Row {
	blankRow := table.Row{
		"",
		"0",
		"0",
		"00:00",
		"0",
		"0",
		"0",
		"0",
		"0",
		"00:00",
	}

	return blankRow
}

func generateBottomRows(monthlyTotals *data.MonthStats) []table.Row {
	var rows []table.Row
	blankRow := table.Row{
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	}
	rows = append(rows, blankRow)
	totalOvertime := helpers.DecimalToTime(monthlyTotals.TotalOvertime.Float64)
	workedTime := strings.Split(monthlyTotals.WorkedTime, ":")
	totalsRow := table.Row{
		"Total",
		strconv.Itoa(monthlyTotals.TotalWeekDays),
		strconv.Itoa(monthlyTotals.WorkedDays),
		workedTime[0] + ":" + workedTime[1],
		strconv.Itoa(monthlyTotals.VacationDays),
		strconv.Itoa(monthlyTotals.SickDays),
		strconv.Itoa(monthlyTotals.WeekendDays),
		strconv.Itoa(monthlyTotals.OffDays),
		strconv.Itoa(monthlyTotals.OverTimeDays),
		totalOvertime,
	}
	rows = append(rows, totalsRow)
	return rows
}

