package registration

import (
	"fTime/data"
	"fmt"
)

func RegisterOvertime(otArgument string, state *data.ReportState, userConfig *data.UserConfig) error {
	if !state.SelectedRecord.EndTime.Valid {
		fmt.Println("end time must be registered before the day is flagged for overtime")
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
	err := RegisterEnd(endTime, state, userConfig)
	if err != nil {
		return err
	}

	return nil
}
