package registration

import (
	"fTime/data"
	"fmt"
	"strconv"
)

func RegisterLunch(lunchString string, state *data.ReportState) error {
	lunchDurationInt, err := strconv.Atoi(lunchString)
	if err != nil {
		return fmt.Errorf("failed to convert lunch duration to numeric value.\nInput format must be <INT(minutes)>")
	}
	lunchDuration := int16(lunchDurationInt)
	if lunchDuration < 1 || lunchDuration > 999 {
		return fmt.Errorf("lunch duration must be between 1 and 99 in format <INT(minutes)>")
	}

	err = data.UpdateLunch(state.SelectedDate, lunchDuration)
	if err != nil {
		return err
	}

	state.SelectedRecord.LunchDuration.Int16 = lunchDuration

	if state.SelectedRecord.EndTime.Valid {
		err = RegisterTotals(state)
		if err != nil {
			return err
		}
		err = rebalanceSucceedingDates(state)
	}

	return nil
}
