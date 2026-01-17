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

	registeredTime, err := helpers.ParseTimeObject(startTime)
	if err != nil {
		return fmt.Errorf("failed to parse start time: \n%v", err)
	}

	if state.SelectedRecord.StartTime.Valid {
		err := data.UpdateStart(
			state.SelectedRecord.WorkDate,
			registeredTime.Format(utils.TimeLayout),
		)
		if err != nil {
			return err
		}
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
	state.SelectedRecord.StartTime.Valid = true
	state.SelectedRecord.DayLength.String = userConfig.DefaultDayLength.String
	state.SelectedRecord.DayLength.Valid = true
	state.MaxDate = state.SelectedDate

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
