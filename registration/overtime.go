package registration

import (
	"fTime/data"
	"fTime/utils"
	"fmt"
	"time"
)

func RegisterOvertime(otArgument string, state *data.ReportState, userConfig *data.UserConfig) error {
	if !state.SelectedRecord.EndTime.Valid {
		fmt.Println("end time must be registered before the day is flagged for overtime")
	}

	registeredDayTotal, err := time.Parse(utils.NonsenseDate, state.SelectedRecord.DayTotal.String)
	if err != nil {
		return fmt.Errorf("failed to parse day total: %v", err)
	}
	strippedOfEightHours := registeredDayTotal.Add(-8 * time.Hour)

	referenceTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if strippedOfEightHours.Before(referenceTime) || strippedOfEightHours.Equal(referenceTime) {
		return fmt.Errorf("day total must be greater than 08:00 hours to register overtime")
	}

	if otArgument == "ot" {
		if state.SelectedRecord.Overtime.Bool {
			return fmt.Errorf("date is already registered as overtime")
		} else {
			state.SelectedRecord.Overtime.Bool = true
		}
	}

	if otArgument == "-ot" {
		if !state.SelectedRecord.Overtime.Bool {
			return fmt.Errorf("date doesn't have an overtime flag to remove")
		} else {
			state.SelectedRecord.Overtime.Bool = false
		}
	}

	endTime := state.SelectedRecord.EndTime.String[:2] + state.SelectedRecord.EndTime.String[3:5]
	err = RegisterEnd(endTime, state, userConfig)
	if err != nil {
		return err
	}

	return nil
}
