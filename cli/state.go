package cli

import (
	"fTime/helpers"
	"fTime/logic"
)

type reportState struct {
	reportUpToDate   bool
	maxDate          string
	maxCompletedDate string
	selectedDate     string
	selectedRecord   helpers.WorkDateRecord
}

func (rs *reportState) GetReportUpToDate() bool {
	return rs.reportUpToDate
}

func (rs *reportState) GetMaxDate() string { return rs.maxDate }

func (rs *reportState) GetMaxCompletedDate() string { return rs.maxCompletedDate }

func (rs *reportState) GetSelectedDate() string {
	return rs.selectedDate
}

func (rs *reportState) GetSelectedRecord() helpers.WorkDateRecord {
	return rs.selectedRecord
}

func setNewState(selectedDate string) (reportState, error) {
	maxCompletedDate, maxDate, err := helpers.GetMaxDates()
	if err != nil {
		return reportState{}, err
	}

	previousCompleted, err := logic.CheckPreviousCompletion(selectedDate, maxCompletedDate)
	if err != nil {
		return reportState{}, err
	}

	recordExists, err := logic.CheckIfDateExists(selectedDate, maxDate)
	if err != nil {
		return reportState{}, err
	}

	var selectedDateRecord helpers.WorkDateRecord
	if recordExists {
		selectedDateRecord, err = helpers.GetOneWorkDateRecord(selectedDate)
		if err != nil {
			return reportState{}, err
		}
	} else {
		selectedDateRecord = helpers.WorkDateRecord{
			WorkDate: selectedDate,
		}
	}

	currentState := reportState{
		reportUpToDate:   previousCompleted,
		maxDate:          maxDate,
		maxCompletedDate: maxCompletedDate,
		selectedDate:     selectedDate,
		selectedRecord:   selectedDateRecord,
	}

	return currentState, nil
}
