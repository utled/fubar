package cli

import (
	"fTime/data"
	"fTime/helpers"
	"fmt"
)

/*func setNewState(selectedDate string, userConfig *data.UserConfig) (data.ReportState, error) {
	maxCompletedDate, maxDate, err := data.GetMaxDates()
	if err != nil {
		return data.ReportState{}, err
	}

	previousCompleted, err := helpers.CheckPreviousCompletion(selectedDate, maxCompletedDate)
	if err != nil {
		return data.ReportState{}, err
	}

	recordExists, err := helpers.CheckIfDateExists(selectedDate, maxDate)
	if err != nil {
		return data.ReportState{}, err
	}

	var selectedDateRecord data.WorkDateRecord
	if recordExists {
		selectedDateRecord, err = data.GetOneWorkDateRecord(selectedDate)
		if err != nil {
			return data.ReportState{}, err
		}
	} else {
		selectedDateRecord = data.WorkDateRecord{
			WorkDate: selectedDate,
		}
	}

	projectedEnd := helpers.CalcProjectedEnd(&selectedDateRecord, userConfig)

	currentState := data.ReportState{
		ReportUpToDate:   previousCompleted,
		MaxDate:          maxDate,
		MaxCompletedDate: maxCompletedDate,
		SelectedDate:     selectedDate,
		SelectedRecord:   &selectedDateRecord,
		ProjectedEnd:     projectedEnd,
	}

	return currentState, nil
}*/

func setNewState(selectedDate string, currentState *data.ReportState, userConfig *data.UserConfig) {
	maxCompletedDate, maxDate, err := data.GetMaxDates()
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

	projectedEnd := helpers.CalcProjectedEnd(&selectedDateRecord, userConfig)
	currentState.ReportUpToDate = previousCompleted
	currentState.MaxDate = maxDate
	currentState.MaxCompletedDate = maxCompletedDate
	currentState.SelectedDate = selectedDate
	currentState.SelectedRecord = &selectedDateRecord
	currentState.ProjectedEnd = projectedEnd

	helpers.PrintSelectedDate(currentState)
}
