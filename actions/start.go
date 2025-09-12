package actions

import (
	"fTime/helpers"
	"fTime/utils"
	"fmt"
	"time"
)

func RegisterStart(startTime string, state *helpers.ReportState) error {
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
		err = helpers.UpdateStart(state.SelectedRecord.WorkDate, registeredTime.Format(utils.TimeLayout))
	} else {
		err = helpers.WriteStart(state.SelectedRecord.WorkDate, registeredTime.Format(utils.TimeLayout))
		if err != nil {
			return err
		}
		state.SelectedRecord.StartTime.String = registeredTime.Format(utils.TimeLayout)
		state.SelectedRecord.StartTime.Valid = true
	}

	if !state.SelectedRecord.EndTime.Valid {
		return nil
	}

	dayTotal, err := helpers.CalcDayTotal(state.SelectedRecord)
	if err != nil {
		return err
	}
	state.SelectedRecord.DayTotal.String = dayTotal
	state.SelectedRecord.DayTotal.Valid = true

	dayBalance, err := helpers.CalcDayBalance(state.SelectedRecord)
	if err != nil {
		return err
	}
	state.SelectedRecord.DayBalance.Float64 = dayBalance
	state.SelectedRecord.DayBalance.Valid = true

	selectedDate, err := time.Parse(utils.DateLayout, state.SelectedDate)
	if err != nil {
		return fmt.Errorf("failed to parse selectedDate.%v", err)
	}
	previousBalance, err := helpers.GetPreviousBalance(selectedDate)
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
	state.SelectedRecord.MovingBalance.Valid = true

	err = helpers.WriteNewBalance(selectedDate.Format(utils.DateLayout), dayTotal, dayBalance, totalBalance)
	if err != nil {
		return err
	}

	maxDateString := state.MaxDate
	maxDate, err := time.Parse(utils.DateLayout, maxDateString)
	if err != nil {
		return fmt.Errorf("failed to parse maxDate.%v", err)
	}

	if selectedDate.Before(maxDate) {
		rebalanceSucceedingDates()
	}

	return nil
}
