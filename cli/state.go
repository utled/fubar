package cli

import (
	"fTime/helpers"
	"fTime/logic"
)

type reportState struct {
	reportUpToDate   bool
	maxCompletedDate string
	selectedDate     string
	selectedRecord   helpers.WorkDateRecord
}

func (ws *reportState) GetReportUpToDate() bool {
	return ws.reportUpToDate
}

func (ws *reportState) GetMaxCompletedDate() string {
	return ws.maxCompletedDate
}

func (ws *reportState) GetSelectedDate() string {
	return ws.selectedDate
}

func (ws *reportState) GetSelectedRecord() helpers.WorkDateRecord {
	return ws.selectedRecord
}

func setNewState(selectedDate string) (reportState, error) {
	recordExists, err := logic.CheckIfDateExists(selectedDate)
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

	previousCompleted, maxCompletedDate, err := logic.CheckPreviousCompletion(selectedDate)
	if err != nil {
		return reportState{}, err
	}

	currentState := reportState{
		reportUpToDate:   previousCompleted,
		maxCompletedDate: maxCompletedDate,
		selectedDate:     selectedDate,
		selectedRecord:   selectedDateRecord,
	}

	return currentState, nil
}
