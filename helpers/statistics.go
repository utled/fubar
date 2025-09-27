package helpers

import (
	"fTime/data"
	"fTime/utils"
	"fmt"
	"strconv"
	"time"
)

func DisplayAllStatistics(state *data.ReportState) error {
	startDate := time.Date(1901, 1, 1, 0, 0, 0, 0, time.UTC).Format(utils.DateLayout)
	endDate := time.Now().Format(utils.DateLayout)

	fullStatistics, err := data.GetFullStatistics(startDate, endDate)
	if err != nil {
		return err
	}

	title := "Statistics for all recorded dates"

	PrintFullStatistics(fullStatistics, title, state)

	return nil
}

func DisplaySumStatistics(year string, state *data.ReportState) error {
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return fmt.Errorf("Invalid argument.\nExpects: 'stats[st] sum[s] <YYYY>'")
	}

	monthlySummary, err := data.GetMonthlySummary(yearInt)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("Monthly Summary for %d", yearInt)

	PrintMonthlySummary(monthlySummary, title, state)

	return nil
}

func DisplayYearStatistics(year string, state *data.ReportState) error {
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return fmt.Errorf("Invalid argument.\nExpects: 'stats[st] year[y] <YYYY>'")
	}

	startDate := time.Date(yearInt, 1, 1, 0, 0, 0, 0, time.UTC).Format(utils.DateLayout)
	endDate := time.Date(yearInt, 12, 31, 0, 0, 0, 0, time.UTC).Format(utils.DateLayout)

	fullStatistics, err := data.GetFullStatistics(startDate, endDate)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("Statistics for %d", yearInt)

	PrintFullStatistics(fullStatistics, title, state)

	return nil
}

func DisplayMonthStatistics(month string, year string, state *data.ReportState) error {
	invalidArgumentMsg := fmt.Errorf("Invalid argument.\nExpects: 'stats[st] month[m] <INT(monthnum)> <YYYY>'")

	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return invalidArgumentMsg
	}

	if monthInt < 1 || monthInt > 12 {
		return fmt.Errorf("Invalid month number.\n Month must be between 1 and 12")
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return invalidArgumentMsg
	}

	startDate := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC).Format(utils.DateLayout)
	firstOfNextMonth := time.Date(yearInt, time.Month(monthInt)+1, 1, 0, 0, 0, 0, time.UTC)
	endDate := firstOfNextMonth.Add(-time.Nanosecond).Format(utils.DateLayout)

	fullStatistics, err := data.GetFullStatistics(startDate, endDate)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("Statistics for %s %s", time.Month(monthInt).String(), year)

	PrintFullStatistics(fullStatistics, title, state)

	return nil
}

func DisplayDaysStatistics(numOfDays string, state *data.ReportState) error {
	numOfDaysInt, err := strconv.Atoi(numOfDays)
	if err != nil {
		return fmt.Errorf("Invalid argument.\nExpects: 'stats[st] day[d] <INT(days)>'")
	}

	today := time.Now()
	startDate := today.AddDate(0, 0, -numOfDaysInt).Format(utils.DateLayout)
	endDate := today.Format(utils.DateLayout)

	fullStatistics, err := data.GetFullStatistics(startDate, endDate)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("Statistics for the last %d days (%s - %s)", numOfDaysInt, startDate, endDate)

	PrintFullStatistics(fullStatistics, title, state)

	return nil
}

func DisplayRangeStatistics(startDateString string, endDateString string, state *data.ReportState) error {
	startDate, err := FormatValidDateString(startDateString)
	if err != nil {
		return err
	}

	endDate, err := FormatValidDateString(endDateString)
	if err != nil {
		return err
	}

	startIsBeforeEnd, err := CheckDateBefore(startDate, endDate)
	if err != nil {
		return err
	}
	if !startIsBeforeEnd {
		return fmt.Errorf("start date must be before end date")
	}

	fullStatistics, err := data.GetFullStatistics(startDate, endDate)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("Statistics for date range %s - %s", startDate, endDate)

	PrintFullStatistics(fullStatistics, title, state)

	return nil
}
