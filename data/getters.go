package data

import (
	"database/sql"
	"fTime/db"
	"fTime/utils"
	"fmt"
	"time"
)

func GetTimesheetRange() error {
	con, err := openDBConnection()
	if err != nil {
		return err
	}

	query := "SELECT * FROM timesheet WHERE workdate between ? AND ?;"
	startDate := "2024-02-01"
	endDate := "2024-02-01"
	response, err := con.Query(query, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
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
			&workDateRecord.DayLength,
			&workDateRecord.DayType,
		)
		if err != nil {
			return fmt.Errorf("failed to serialize range of records to struct: %v", err)
		}
		timesheet = append(timesheet, *workDateRecord)
	}
	/*	for idx := len(timesheet) - 1; idx >= 0; idx-- {
		fmt.Println(timesheet[idx])
	}*/

	return nil
}

func GetOneWorkDateRecord(queryDate string) (record WorkDateRecord, err error) {
	con, err := openDBConnection()
	if err != nil {
		return WorkDateRecord{}, err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	query := "SELECT * FROM timesheet WHERE workdate = ?"
	response := con.QueryRow(query, queryDate)

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
		&workDateRecord.DayLength,
		&workDateRecord.DayType,
	)
	if err != nil {
		return WorkDateRecord{}, fmt.Errorf("failed to serialize record to struct%v", err)
	}

	return *workDateRecord, nil

}

func GetMaxDates() (maxCompletedDate string, maxDate string, err error) {
	con, err := openDBConnection()
	if err != nil {
		return "", "", err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	query := "SELECT MAX(workdate) FROM timesheet WHERE end_time IS NOT NULL;"
	response := con.QueryRow(query)
	err = response.Scan(&maxCompletedDate)
	if err != nil {
		return "", "", fmt.Errorf("failed to read max completed date: %v", err)
	}

	query = "SELECT MAX(workdate) FROM timesheet;"
	response = con.QueryRow(query)
	err = response.Scan(&maxDate)
	if err != nil {
		return "", "", fmt.Errorf("failed to read max date: %v", err)
	}

	return maxCompletedDate, maxDate, nil
}

func GetUserConfig() (config UserConfig, err error) {
	con, err := openDBConnection()
	if err != nil {
		return UserConfig{}, err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	query := "SELECT * FROM userconfig limit 1;"
	response := con.QueryRow(query)

	userConfig := &UserConfig{}

	err = response.Scan(
		&userConfig.DefaultLunch,
		&userConfig.DefaultDayLength,
		&userConfig.OffStart,
		&userConfig.OffEnd,
		&userConfig.OffType,
	)
	if err != nil {
		return UserConfig{}, fmt.Errorf("failed to serialize user config to struct%v", err)
	}

	return *userConfig, nil
}

func GetPreviousBalance(selectedDate time.Time) (previousBalance float64, err error) {
	con, err := openDBConnection()
	if err != nil {
		return 0.0, err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	previousDate := selectedDate.AddDate(0, 0, -1)
	previousDateString := previousDate.Format(utils.DateLayout)
	query := "SELECT moving_balance FROM timesheet WHERE workdate = ?;"
	response := con.QueryRow(query, previousDateString)
	err = response.Scan(&previousBalance)
	if err != nil {
		return 0.0, fmt.Errorf("failed to read moving balance: %v", err)
	}

	return previousBalance, nil
}
