package registration

import (
	"fTime/data"
	"fTime/helpers"
	"fTime/utils"
	"fmt"
	"time"
)

func RegisterTotals(state *data.ReportState) error {
	dayTotal, err := helpers.CalcDayTotal(state.SelectedRecord)
	if err != nil {
		return err
	}
	state.SelectedRecord.DayTotal.String = dayTotal

	dayBalance, err := helpers.CalcDayBalance(state.SelectedRecord)
	if err != nil {
		return err
	}
	state.SelectedRecord.DayBalance.Float64 = dayBalance

	selectedDate, err := time.Parse(utils.DateLayout, state.SelectedDate)
	if err != nil {
		return fmt.Errorf("failed to parse selectedDate.%v", err)
	}
	previousBalance, err := data.GetPreviousBalance(selectedDate)
	if err != nil {
		return err
	}
	var totalBalance float64
	if state.SelectedRecord.Overtime.Bool {
		totalBalance = previousBalance
	} else {
		totalBalance = helpers.CalcTotalBalance(state.SelectedRecord, previousBalance)
	}
	state.SelectedRecord.MovingBalance.Float64 = totalBalance

	err = data.WriteNewBalance(selectedDate.Format(utils.DateLayout), dayTotal, dayBalance, totalBalance)
	if err != nil {
		return err
	}

	return nil
}
