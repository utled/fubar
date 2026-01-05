package helpers

import (
	"fmt"
	"fubar/data"
	"fubar/utils"
	"strconv"
	"time"
)

func SetDateRangeFromDates(startDate string, endDate string) (dateRange []*data.WorkDateRecord, err error) {
	formattedStart, err := FormatValidDateString(startDate)
	if err != nil {
		return []*data.WorkDateRecord{}, err
	}
	formattedEnd, err := FormatValidDateString(endDate)
	if err != nil {
		return []*data.WorkDateRecord{}, err
	}

	startIsBeforeEnd, err := CheckDateBefore(formattedStart, formattedEnd)
	if err != nil {
		return []*data.WorkDateRecord{}, err
	}

	if !startIsBeforeEnd {
		return []*data.WorkDateRecord{}, fmt.Errorf("start date must be before end date:%v", err)
	}

	dateRange, err = data.GetTimesheetRange(formattedStart, formattedEnd)
	if err != nil {
		return []*data.WorkDateRecord{}, err
	}

	return dateRange, nil
}

func SetDateRangeFromDayCount(days string) (dateRange []*data.WorkDateRecord, err error) {
	dayCount, err := strconv.Atoi(days)
	if err != nil {
		return []*data.WorkDateRecord{}, fmt.Errorf("input includes non-numerical characters")
	}

	if dayCount < 1 {
		return []*data.WorkDateRecord{}, fmt.Errorf("days must be greater than zero")
	}

	today := time.Now()
	startDate := today.AddDate(0, 0, -dayCount)

	dateRange, err = data.GetTimesheetRange(startDate.Format(utils.DateLayout), today.Format(utils.DateLayout))
	if err != nil {
		return []*data.WorkDateRecord{}, err
	}

	return dateRange, nil
}
