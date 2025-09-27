package registration

import (
	"database/sql"
	"fTime/data"
	"fTime/utils"
	"fmt"
	"time"
)

func makeOffDay(dateString string, previousBalance float64, offType string) data.WorkDateRecord {
	dateRecord := data.WorkDateRecord{
		WorkDate:       dateString,
		StartTime:      sql.NullString{String: "00:00:00", Valid: true},
		EndTime:        sql.NullString{String: "00:00:00", Valid: true},
		LunchDuration:  sql.NullInt16{Int16: 0, Valid: true},
		DayTotal:       sql.NullString{String: "00:00:00", Valid: true},
		DayBalance:     sql.NullFloat64{Float64: 0.0, Valid: true},
		Overtime:       sql.NullBool{Bool: false, Valid: true},
		TotalBalance:   sql.NullFloat64{Float64: previousBalance, Valid: true},
		AdditionalTime: sql.NullInt16{Int16: 0, Valid: true},
		DayLength:      sql.NullString{String: "00:00:00", Valid: true},
		DayType:        sql.NullString{String: offType, Valid: true},
	}

	return dateRecord
}

func makeFullDay(dateString string, previousBalance float64, offType string) data.WorkDateRecord {
	dateRecord := data.WorkDateRecord{
		WorkDate:       dateString,
		StartTime:      sql.NullString{String: "08:00:00", Valid: true},
		EndTime:        sql.NullString{String: "16:40:00", Valid: true},
		LunchDuration:  sql.NullInt16{Int16: 40, Valid: true},
		DayTotal:       sql.NullString{String: "08:00:00", Valid: true},
		DayBalance:     sql.NullFloat64{Float64: 0.0, Valid: true},
		Overtime:       sql.NullBool{Bool: false, Valid: true},
		TotalBalance:   sql.NullFloat64{Float64: previousBalance, Valid: true},
		AdditionalTime: sql.NullInt16{Int16: 0, Valid: true},
		DayLength:      sql.NullString{String: "08:00:00", Valid: true},
		DayType:        sql.NullString{String: offType, Valid: true},
	}

	return dateRecord
}

func RegisterBackflush(dayType string, state *data.ReportState, userConfig *data.UserConfig) error {
	lastCompletedDate, err := time.Parse(utils.DateLayout, state.MaxCompletedDate)
	if err != nil {
		return fmt.Errorf("failed to parse last completed date: %v", err)
	}
	maxDate, err := time.Parse(utils.DateLayout, state.MaxDate)
	if err != nil {
		return fmt.Errorf("failed to parse max date: %v", err)
	}

	if lastCompletedDate.Before(maxDate) {
		err = data.DeleteRecord(maxDate.Format(utils.DateLayout))
		if err != nil {
			return fmt.Errorf("failed to delete last incompleted date: %v", err)
		}
	}

	today := time.Now()

	previousBalance, err := data.GetPreviousBalance(lastCompletedDate.AddDate(0, 0, 1))
	if err != nil {
		return err
	}

	backflushRange := make([]data.WorkDateRecord, 0)

	for currentDate := lastCompletedDate.AddDate(0, 0, 1); currentDate.Before(today); currentDate = currentDate.AddDate(0, 0, 1) {
		weekday := currentDate.Weekday()
		dateString := currentDate.Format(utils.DateLayout)
		if weekday == time.Saturday || weekday == time.Sunday {
			backflushRange = append(backflushRange, makeOffDay(dateString, previousBalance, "wknd"))
		} else {
			if dayType == "norm" {
				backflushRange = append(backflushRange, makeFullDay(dateString, previousBalance, "norm"))
			} else {
				backflushRange = append(backflushRange, makeOffDay(dateString, previousBalance, dayType))
			}
		}
	}

	err = data.WriteBackflush(&backflushRange)
	if err != nil {
		return err
	}

	nextDay := today.AddDate(0, 0, 1)
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

	return nil
}
