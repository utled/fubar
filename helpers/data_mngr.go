package helpers

import (
	"database/sql"
	"fTime/db"
	"fmt"
)

type StatusProvider interface {
	GetReportUpToDate() bool
	GetMaxDate() string
	GetMaxCompletedDate() string
	GetSelectedDate() string
	GetSelectedRecord() WorkDateRecord
}

type WorkDateRecord struct {
	WorkDate       string
	StartTime      sql.NullString
	EndTime        sql.NullString
	LunchDuration  sql.NullInt16
	DayTotal       sql.NullString
	DayBalance     sql.NullFloat64
	Overtime       sql.NullBool
	MovingBalance  sql.NullFloat64
	AdditionalTime sql.NullInt16
	SickDay        sql.NullBool
	DayLength      sql.NullString
}

func GetTimesheet() {
	query := "SELECT * FROM timesheet WHERE workdate between ? AND ?;"
	startDate := "2024-02-01"
	endDate := "2024-02-01"
	response, err := db.GetMultiRecords(query, startDate, endDate)
	if err != nil {
		fmt.Println(err)
	}

	var timesheet []WorkDateRecord

	for response.Next() {
		workDateRecord := &WorkDateRecord{}
		err := response.Scan(
			&workDateRecord.WorkDate,
			&workDateRecord.StartTime,
			&workDateRecord.EndTime,
			&workDateRecord.LunchDuration,
			&workDateRecord.DayTotal,
			&workDateRecord.DayBalance,
			&workDateRecord.Overtime,
			&workDateRecord.MovingBalance,
			&workDateRecord.AdditionalTime,
			&workDateRecord.SickDay,
			&workDateRecord.DayLength,
		)
		if err != nil {
			fmt.Println(err)
		}
		timesheet = append(timesheet, *workDateRecord)
	}
	/*	for idx := len(timesheet) - 1; idx >= 0; idx-- {
		fmt.Println(timesheet[idx])
	}*/
}

func GetOneWorkDateRecord(queryDate string) (workdateRecord WorkDateRecord, err error) {
	query := "SELECT * FROM timesheet WHERE workdate = ?"

	response, err := db.GetOneRecord(query, queryDate)
	if err != nil {
		fmt.Println(err)
	}

	workDateRecord := &WorkDateRecord{}

	err = response.Scan(
		&workDateRecord.WorkDate,
		&workDateRecord.StartTime,
		&workDateRecord.EndTime,
		&workDateRecord.LunchDuration,
		&workDateRecord.DayTotal,
		&workDateRecord.DayBalance,
		&workDateRecord.Overtime,
		&workDateRecord.MovingBalance,
		&workDateRecord.AdditionalTime,
		&workDateRecord.SickDay,
		&workDateRecord.DayLength,
	)
	if err != nil {
		return *workDateRecord, fmt.Errorf("failed to serialize query response to struct%v", err)
	}

	/*fmt.Println(workDateRecord)
	fmt.Println(workDateRecord.MovingBalance.Float64)*/

	return *workDateRecord, nil

}

func GetMaxDates() (maxCompletedDate string, maxDate string, err error) {
	query := "SELECT MAX(workdate) FROM timesheet WHERE end_time IS NOT NULL;"
	maxCompletedDate, err = db.GetOneValue(query, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to query maxCompletedDate %v", err)
	}

	query = "SELECT MAX(workdate) FROM timesheet;"
	maxDate, err = db.GetOneValue(query, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to query maxDate %v", err)
	}

	return maxCompletedDate, maxDate, nil
}
