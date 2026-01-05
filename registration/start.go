package registration

import (
	"fmt"
	"fubar/data"
	"fubar/helpers"
	"fubar/utils"
)

func RegisterStart(startTime string, state *data.ReportState, userConfig *data.UserConfig) error {
	if !state.ReportUpToDate {
		return fmt.Errorf("can't start selected date.\nAll previous dates must be up to date.")
	}

	formattedTime, err := helpers.FormatValidTimeString(startTime)
	if err != nil {
		return fmt.Errorf("failed to format start time.%v", err)
	}
	registeredTime, err := helpers.ParseTimeObject(formattedTime)
	if err != nil {
		return fmt.Errorf("failed to parse start time.%v", err)
	}

	if state.SelectedRecord.StartTime.Valid {
		err = data.UpdateStart(
			state.SelectedRecord.WorkDate,
			registeredTime.Format(utils.TimeLayout),
		)
	} else {
		err = data.WriteStart(
			state.SelectedRecord.WorkDate,
			registeredTime.Format(utils.TimeLayout),
			userConfig.DefaultDayLength.String,
		)
		if err != nil {
			return err
		}
	}
	state.SelectedRecord.StartTime.String = registeredTime.Format(utils.TimeLayout)

	if !state.SelectedRecord.EndTime.Valid {
		return nil
	}

	err = RegisterTotals(state)
	if err != nil {
		return err
	}

	selectedBeforeMax, err := helpers.CheckDateBefore(state.SelectedDate, state.MaxCompletedDate)
	if err != nil {
		return err
	}

	if selectedBeforeMax {
		err = rebalanceSucceedingDates(state)
		if err != nil {
			return err
		}
	}

	return nil
}
