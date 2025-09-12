package cli

import (
	"fTime/actions"
	"fTime/helpers"
)

func setNewState(selectedDate string, userConfig *helpers.UserConfig) (helpers.ReportState, error) {
	maxCompletedDate, maxDate, err := helpers.GetMaxDates()
	if err != nil {
		return helpers.ReportState{}, err
	}

	previousCompleted, err := actions.CheckPreviousCompletion(selectedDate, maxCompletedDate)
	if err != nil {
		return helpers.ReportState{}, err
	}

	recordExists, err := actions.CheckIfDateExists(selectedDate, maxDate)
	if err != nil {
		return helpers.ReportState{}, err
	}

	var selectedDateRecord helpers.WorkDateRecord
	if recordExists {
		selectedDateRecord, err = helpers.GetOneWorkDateRecord(selectedDate)
		if err != nil {
			return helpers.ReportState{}, err
		}
	} else {
		selectedDateRecord = helpers.WorkDateRecord{
			WorkDate: selectedDate,
		}
	}

	projectedEnd := helpers.CalcProjectedEnd(&selectedDateRecord, userConfig)

	currentState := helpers.ReportState{
		ReportUpToDate:   previousCompleted,
		MaxDate:          maxDate,
		MaxCompletedDate: maxCompletedDate,
		SelectedDate:     selectedDate,
		SelectedRecord:   &selectedDateRecord,
		ProjectedEnd:     projectedEnd,
	}

	return currentState, nil
}
