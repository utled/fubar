package logic

import (
	"fTime/helpers"
	"fTime/utils"
	"fmt"
	"time"
)

func RegisterStart(startTime string, status helpers.StatusProvider) error {
	if !status.GetReportUpToDate() {
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

	dateRecord := status.GetSelectedRecord()
	dateRecord.StartTime.String = registeredTime.Format(utils.TimeLayout)
	dateRecord.StartTime.Valid = true

	if !dateRecord.EndTime.Valid {
		query := "UPDATE timesheet SET start_time = ? WHERE workdate = ?"
		fmt.Println(query)
		//Execute DB update
		return nil
	}

	dayTotal, err := helpers.CalcDayTotal(&dateRecord)
	if err != nil {
		return err
	}
	dateRecord.DayTotal.String = dayTotal
	dateRecord.DayTotal.Valid = true

	dayBalance, err := helpers.CalcDayBalance(&dateRecord)
	if err != nil {
		return err
	}
	dateRecord.DayBalance.Float64 = dayBalance
	dateRecord.DayBalance.Valid = true

	selectedDate, err := time.Parse(utils.DateLayout, status.GetSelectedDate())
	if err != nil {
		return fmt.Errorf("failed to parse selectedDate.%v", err)
	}
	previousBalance, err := helpers.GetPreviousBalance(selectedDate)
	if err != nil {
		return err
	}

	var totalBalance float64
	if dateRecord.Overtime.Bool {
		totalBalance = previousBalance
	} else {
		totalBalance = helpers.CalcTotalBalance(&dateRecord, previousBalance)
	}
	dateRecord.MovingBalance.Float64 = totalBalance
	dateRecord.MovingBalance.Valid = true

	err = helpers.WriteNewBalance(selectedDate.Format(utils.DateLayout), dayTotal, dayBalance, totalBalance)
	if err != nil {
		return err
	}

	maxDateString := status.GetMaxDate()
	maxDate, err := time.Parse(utils.DateLayout, maxDateString)
	if err != nil {
		return fmt.Errorf("failed to parse maxDate.%v", err)
	}
	if selectedDate.Before(maxDate) {
		rebalanceSucceedingDates()
	}

	return nil
}
