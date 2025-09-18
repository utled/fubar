package registration

import (
	"fTime/data"
	"fTime/utils"
	"fmt"
	"time"
)

func rebalanceSucceedingDates(state *data.ReportState) error {
	selectedDate, err := time.Parse(utils.DateLayout, state.SelectedDate)
	if err != nil {
		return fmt.Errorf("failed to parse selected date: %v", err)
	}

	maxCompletedDate, err := time.Parse(utils.DateLayout, state.MaxCompletedDate)
	if err != nil {
		return fmt.Errorf("failed to parse max completed date: %v", err)
	}

	var succeedingDates []string

	for currentDate := selectedDate.AddDate(0, 0, 1); currentDate.Before(maxCompletedDate) || currentDate.Equal(maxCompletedDate); currentDate = currentDate.AddDate(0, 0, 1) {
		succeedingDates = append(succeedingDates, currentDate.Format(utils.DateLayout))
	}

	err = data.UpdateTotalBalance(&succeedingDates, state.SelectedRecord.TotalBalance.Float64)
	if err != nil {
		return err
	}

	return nil
}
