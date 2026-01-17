package registration

import (
	"fmt"
	"fubar/data"
)

func RegisterAdditionalTime(additionalTime int, state *data.ReportState) error {
	if !state.SelectedRecord.EndTime.Valid {
		return fmt.Errorf("end time must be registered before additional time")
	}

	err := data.UpdateAdditionalTime(state.SelectedDate, additionalTime)
	if err != nil {
		return err
	}

	state.SelectedRecord.AdditionalTime.Int16 = int16(additionalTime)
	state.SelectedRecord.AdditionalTime.Valid = true

	err = RegisterTotals(state)
	if err != nil {
		return err
	}
	err = rebalanceSucceedingDates(state)
	if err != nil {
		return err
	}

	return nil
}
