package helpers

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

	dateDiff := today.Sub(maxCompletedDate)
	if dateDiff.Hours() > 24 {
		isCompleted = false
	} else {
		isCompleted = true
	}

	return isCompleted, nil
}

func CheckIfDateExists(dateString string, maxDateString string) (dateExists bool, err error) {
	dateToCheck, err := time.Parse(utils.DateLayout, dateString)
	if err != nil {
		return false, fmt.Errorf("failed to parse selected date %v", err)
	}
	maxDate, err := time.Parse(utils.DateLayout, maxDateString)
	if err != nil {
		return false, fmt.Errorf("failed to parse max date %v", err)
	}
	dateDiff := dateToCheck.Sub(maxDate)

	dateExists = dateDiff.Hours() <= 0

	return dateExists, nil
}

func CheckDateInFuture(dateString string) (dateInFuture bool, err error) {
	dateToCheck, err := time.Parse(utils.DateLayout, dateString)
	if err != nil {
		return false, fmt.Errorf("failed to parse selected date %v", err)
	}
	dateInFuture = time.Now().Before(dateToCheck)
	return dateInFuture, err
}

func CheckDateBefore(startDate string, endDate string) (dateBefore bool, err error) {
	parsedStart, err := time.Parse(utils.DateLayout, startDate)
	if err != nil {
		return false, fmt.Errorf("failed to parse start date %v", err)
	}
	parsedEnd, err := time.Parse(utils.DateLayout, endDate)
	if err != nil {
		return false, fmt.Errorf("failed to parse end date %v", err)
	}
	dateBefore = parsedStart.Before(parsedEnd)

	return dateBefore, err
}
