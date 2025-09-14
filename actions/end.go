package actions

import (
	"fTime/helpers"
	"fTime/utils"
	"fmt"
	"time"
)

func RegisterEnd(endTime string, state *helpers.ReportState, userConfig *helpers.UserConfig) error {
	if !state.SelectedRecord.StartTime.Valid {
		return fmt.Errorf("can't end selected date.\nstart time must be registered first.")
	}

	formattedTime, err := helpers.FormatValidTimeString(endTime)
	if err != nil {
		return fmt.Errorf("failed to format end time.%v", err)
	}
	registeredTime, err := helpers.ParseTimeObject(formattedTime)
	if err != nil {
		return fmt.Errorf("failed to parse end time: %v", err)
	}

	var lunchDuration int16
	if state.SelectedRecord.LunchDuration.Valid {
		lunchDuration = state.SelectedRecord.LunchDuration.Int16
	} else {
		lunchDuration = userConfig.DefaultLunch.Int16
	}

	var additionalTime int16
	if state.SelectedRecord.AdditionalTime.Valid {
		additionalTime = state.SelectedRecord.AdditionalTime.Int16
	} else {
		additionalTime = 0
	}

	err = helpers.UpdateEnd(
		state.SelectedDate,
		registeredTime.Format(utils.TimeLayout),
		state.SelectedRecord.Overtime.Bool,
		lunchDuration,
		additionalTime,
	)
	if err != nil {
		return err
	}
	state.SelectedRecord.EndTime.String = registeredTime.Format(utils.TimeLayout)

	err = RegisterTotals(state)
	if err != nil {
		return err
	}

	if state.SelectedDate == state.MaxDate {
		parsedDate, err := time.Parse(utils.DateLayout, state.SelectedDate)
		if err != nil {
			return fmt.Errorf("failed to parse selected date.%v", err)
		}
		nextDay := parsedDate.AddDate(0, 0, 1)
		if userConfig.OffStart.String != "" {
			parsedScheduledStart, err := time.Parse(utils.DateLayout, userConfig.OffStart.String)
			if err != nil {
				return fmt.Errorf("failed to parse scheduled start date.%v", err)
			}
			if nextDay == parsedScheduledStart {
				err = RegisterOffPeriod(nextDay, userConfig, state)
				if err != nil {
					return err
				}
			}
		}
		if nextDay.Weekday() == time.Saturday {
			err = RegisterWeekend(nextDay, userConfig, state)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
