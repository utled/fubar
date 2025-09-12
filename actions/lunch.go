package actions

import (
	"fTime/helpers"
	"fmt"
	"strconv"
)

func RegisterLunch(lunchString string, state *helpers.ReportState) error {
	lunchDurationInt, err := strconv.Atoi(lunchString)
	if err != nil {
		return fmt.Errorf("failed to convert lunch duration to numeric value.\nInput format must be <MM>")
	}
	lunchDuration := int16(lunchDurationInt)
	if lunchDuration < 1 || lunchDuration > 59 {
		return fmt.Errorf("lunch duration must be between 1 and 59 in format <MM>")
	}

	err = helpers.UpdateLunch(state.SelectedDate, lunchDuration)
	if err != nil {
		return err
	}

	state.SelectedRecord.LunchDuration.Int16 = lunchDuration

	if state.SelectedRecord.EndTime.Valid {
		err = RegisterTotals(state)
		if err != nil {
			return err
		}
	}

	return nil
}
