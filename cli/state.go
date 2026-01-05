package cli

import (
	"fmt"
	"fubar/data"
	"fubar/helpers"
	"fubar/utils"
	"time"
)

func calcProjectedEnd(dateRecord *data.WorkDateRecord, userConfig *data.UserConfig) string {
	startTime, err := helpers.ParseTimeObject(dateRecord.StartTime.String)
	if err != nil {
		return ""
	}

	dayLength, err := helpers.ParseTimeObject(dateRecord.DayLength.String)
	if err != nil {
		return ""
	}

	var lunchDuration time.Time
	if dateRecord.LunchDuration.Valid {
		lunchDuration = dayLength.Add(time.Minute * time.Duration(dateRecord.LunchDuration.Int16))
		fmt.Println(lunchDuration)
	} else {
		lunchDuration = dayLength.Add(time.Minute * time.Duration(userConfig.DefaultLunch.Int16))
	}

	addHour := time.Duration(lunchDuration.Hour()) * time.Hour
	addMinute := time.Duration(lunchDuration.Minute()) * time.Minute
	projectedEnd := startTime.Add(addHour + addMinute)

	return projectedEnd.Format(utils.TimeLayout)
}

func setNewState(selectedDate string, currentState *data.ReportState, userConfig *data.UserConfig) {
	maxCompletedDate, maxDate, err := data.GetMaxDates()
	if err != nil {
		fmt.Println("failed to set new state.\n", err)
	}

	totalBalance, err := data.GetCurrentTotalBalance()
	if err != nil {
		fmt.Println("failed to set new state.\n", err)
	}

	previousCompleted, err := helpers.CheckPreviousCompletion(selectedDate, maxCompletedDate)
	if err != nil {
		fmt.Println("failed to set new state.\n", err)
	}

	recordExists, err := helpers.CheckIfDateExists(selectedDate, maxDate)
	if err != nil {
		fmt.Println("failed to set new state.\n", err)
	}

	var selectedDateRecord data.WorkDateRecord
	if recordExists {
		selectedDateRecord, err = data.GetOneWorkDateRecord(selectedDate)
		if err != nil {
			fmt.Println("failed to set new state.\n", err)
		}
	} else {
		selectedDateRecord = data.WorkDateRecord{
			WorkDate: selectedDate,
		}
	}

	projectedEnd := calcProjectedEnd(&selectedDateRecord, userConfig)
	currentState.ReportUpToDate = previousCompleted
	currentState.MaxDate = maxDate
	currentState.MaxCompletedDate = maxCompletedDate
	currentState.TotalBalance = totalBalance
	currentState.SelectedDate = selectedDate
	currentState.SelectedRecord = &selectedDateRecord
	currentState.ProjectedEnd = projectedEnd

	helpers.PrintSelectedDate(currentState)
}
