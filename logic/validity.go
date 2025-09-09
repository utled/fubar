package logic

import (
	"fTime/utils"
	"fmt"
	"time"
)

func CheckPreviousCompletion(selectedDate string, maxCompletedDateString string) (isCompleted bool, err error) {
	today, err := time.Parse(utils.DateLayout, selectedDate)
	if err != nil {
		return false, fmt.Errorf("failed to parse selected date%v", err)
	}

	maxCompletedDate, err := time.Parse(utils.DateLayout, maxCompletedDateString)
	if err != nil {
		return false, fmt.Errorf("failed to parse max completed date%v", err)
	}

	var previousCompleted bool
	dateDiff := today.Sub(maxCompletedDate)
	if dateDiff.Hours() > 24 {
		previousCompleted = false
	} else {
		previousCompleted = true
	}

	return previousCompleted, nil
}

func CheckIfDateExists(dateString string, maxDateString string) (dateExists bool, err error) {
	dateToCheck, err := time.Parse(utils.DateLayout, dateString)
	if err != nil {
		return false, fmt.Errorf("failed to parse date %v", err)
	}
	maxDate, err := time.Parse("2006-01-02", maxDateString)
	dateDiff := dateToCheck.Sub(maxDate)

	return dateDiff.Hours() < 0, nil
}
