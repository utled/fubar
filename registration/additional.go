package registration

import (
	"fTime/data"
	"fmt"
	"strconv"
)

func RegisterAdditionalTime(additionalTimeString string, state *data.ReportState) error {
	additionalTimeInt, err := strconv.Atoi(additionalTimeString)
	if err != nil {
		return fmt.Errorf("failed to convert input to numeric value.\nInput format must be <MM>")
	}
	additionalTime := int16(additionalTimeInt)
	if additionalTime < 0 {
		return fmt.Errorf("additionalTime can't be a negative value")
	}

	err = data.UpdateAdditionalTime(state.SelectedDate, additionalTime)
	if err != nil {
		return err
	}

	state.SelectedRecord.AdditionalTime.Int16 = additionalTime

	err = RegisterTotals(state)
	if err != nil {
		return err
	}

	return nil
}
