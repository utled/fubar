package tui

import (
	"fmt"
	"fubar/data"
	"fubar/helpers"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

const (
	idxStart int = iota
	idxLunch
	idxEnd
	idxAdditional
	idxDayType
	idxOvertime
	idxDayTotal
	idxDayBalance
)

type sessionState int

const (
	stateDaily sessionState = iota
	stateBackflush
	stateStatistics
	stateSchedule
	stateConfig
	stateConfirm
)

type dayType int

const (
	norm dayType = iota
	wknd
	off
	vac
	sic
)

type confirmationType int

const (
	deleteDate confirmationType = iota
	deleteSchedule
)

type confirmationDetails struct {
	confirmationType confirmationType
	confirmationMsg  string
}

type statsSelection int

const (
	graphDisplay statsSelection = iota
	tableDisplay
)

const (
	idxWorkedDays int = iota
	idxWeekdays
	idxWorkedTime
	idxAvgStart
	idxAvgEnd
	idxAvgLunch
	idxSickDays
	idxVacDays
	idxOTDays
	idxTotalOT
	idxAvgOT
)

type statsDetails struct {
	displayType    statsSelection
	minYear        int
	maxYear        int
	yearSelection  textinput.Model
	graphArea      viewport.Model
	allSumFields   []textinput.Model
	monthSumFields []textinput.Model
	table          table.Model
	tableTotals    data.MonthStats
}

func (d dayType) String() string {
	return [...]string{"norm", "wknd", "off", "vac", "sic"}[d]
}

func (d dayType) Next() dayType {
	return (d + 1) % 5
}

func ParseDayType(s string) dayType {
	switch s {
	case "wknd":
		return wknd
	case "off":
		return off
	case "vac":
		return vac
	case "sic":
		return sic
	default:
		return norm
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
