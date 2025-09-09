package logic

import (
	"fTime/helpers"
	"fmt"
	"time"
)

func CheckPreviousCompletion() (isCompleted bool, maxCompletedString string, err error) {
	todayString := time.Now().Format("2006-01-02")
	today, err := time.Parse("2006-01-02", todayString)
	if err != nil {
		return false, "", fmt.Errorf("failed to parse todays date%v", err)
	}

	maxCompletedString, err = helpers.GetMaxCompletedDate()
	if err != nil {
		return false, "", err
	}
	maxCompleted, err := time.Parse("2006-01-02", maxCompletedString)
	if err != nil {
		return false, "", fmt.Errorf("failed to parse max completed date%v", err)
	}

	var previousCompleted bool
	dateDiff := today.Sub(maxCompleted)
	if dateDiff.Hours() > 24 {
		previousCompleted = false
	} else {
		previousCompleted = true
	}

	return previousCompleted, maxCompletedString, nil
}

func CheckIfDateExists(dateString string) (dateExists bool, err error) {
	dateToCheck, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return false, fmt.Errorf("failed to parse date %v", err)
	}
	maxDateString, err := helpers.GetMaxDate()
	if err != nil {
		fmt.Println("GetMaxDate Error:", err)
	}
	maxDate, err := time.Parse("2006-01-02", maxDateString)
	dateDiff := dateToCheck.Sub(maxDate)

	return dateDiff.Hours() < 0, nil
}
