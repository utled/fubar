package registration

import (
	"fTime/data"
	"fmt"
)

func DeleteDate(state *data.ReportState) error {
	if state.SelectedDate != state.MaxDate {
		return fmt.Errorf("cannot delete record, selected date is not max date")
	}

	err := data.DeleteRecord(state.SelectedDate)
	if err != nil {
		return err
	}

	return nil
}
