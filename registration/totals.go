package registration

import (
	"fmt"
	"fubar/data"
	"fubar/helpers"
	"fubar/utils"
	"time"
)

func calcDayTotal(dateRecord *data.WorkDateRecord) (string, error) {
	startTime, err := helpers.ParseTimeObject(dateRecord.StartTime.String)
	if err != nil {
		return "", fmt.Errorf("failed to parse start time")
	}

	endTime, err := helpers.ParseTimeObject(dateRecord.EndTime.String)
	if err != nil {
		return "", fmt.Errorf("failed to parse end time")
	}

	timeDiff := endTime.Sub(startTime)
	dayTotal := timeDiff -
		(time.Duration(dateRecord.LunchDuration.Int16) * time.Minute) +
		(time.Duration(dateRecord.AdditionalTime.Int16) * time.Minute)

	totalSeconds := int(dayTotal.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	dayTotalString := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	return dayTotalString, nil
}

func calcDayBalance(dateRecord *data.WorkDateRecord) (float64, error) {
	totalTime, err := helpers.ParseTimeObject(dateRecord.DayTotal.String)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse day total time")
	}

	defaultDayLength, err := helpers.ParseTimeObject(dateRecord.DayLength.String)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse day default length")
	}

	timeDiff := totalTime.Sub(defaultDayLength)

	return timeDiff.Hours(), nil
}

func calcTotalBalance(dateRecord *data.WorkDateRecord, previousTotal float64) float64 {
	if dateRecord.Overtime.Bool {
		return previousTotal
	}

	newTotalBalance := dateRecord.DayBalance.Float64 + previousTotal

	return newTotalBalance
}
func RegisterTotals(state *data.ReportState) error {
	dayTotal, err := calcDayTotal(state.SelectedRecord)
	if err != nil {
		return err
	}
	state.SelectedRecord.DayTotal.String = dayTotal
	state.SelectedRecord.DayTotal.Valid = true

	dayBalance, err := calcDayBalance(state.SelectedRecord)
	if err != nil {
		return err
	}
	state.SelectedRecord.DayBalance.Float64 = dayBalance
	state.SelectedRecord.DayBalance.Valid = true

	selectedDate, err := time.Parse(utils.DateLayout, state.SelectedDate)
	if err != nil {
		return fmt.Errorf("failed to parse selectedDate.%v", err)
	}
	previousBalance, err := data.GetPreviousBalance(selectedDate)
	if err != nil {
		return err
	}

	var totalBalance float64
	if state.SelectedRecord.Overtime.Bool || state.SelectedRecord.DayType.String != "norm" {
		totalBalance = previousBalance
	} else {
		totalBalance = calcTotalBalance(state.SelectedRecord, previousBalance)
	}

	err = data.WriteNewBalance(selectedDate.Format(utils.DateLayout), dayTotal, dayBalance, totalBalance)
	if err != nil {
		return err
	}

	state.SelectedRecord.TotalBalance.Float64 = totalBalance
	state.SelectedRecord.TotalBalance.Valid = true

	return nil
}
