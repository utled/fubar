package registration

import (
	"fmt"
	"fubar/data"
)

func RegisterLunch(lunchDuration int, state *data.ReportState) error {
	if !state.SelectedRecord.StartTime.Valid {
		return fmt.Errorf("start time has not been registered")
	}

	if lunchDuration < 1 || lunchDuration > 999 {
		return fmt.Errorf("lunch duration must be between 1 and 99 in format <INT(minutes)>")
	}

	err := data.UpdateLunch(state.SelectedDate, lunchDuration)
	if err != nil {
		return err
	}

	state.SelectedRecord.LunchDuration.Int16 = int16(lunchDuration)
	state.SelectedRecord.LunchDuration.Valid = true

	if state.SelectedRecord.EndTime.Valid {
		err = RegisterTotals(state)
		if err != nil {
			return err
		}
		err = rebalanceSucceedingDates(state)
	}

	return nil
}
