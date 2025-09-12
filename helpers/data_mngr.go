package helpers

import (
	"database/sql"
	"fTime/db"
	"fTime/utils"
	"fmt"
	"time"
)

type ReportState struct {
	ReportUpToDate   bool
	MaxDate          string
	MaxCompletedDate string
	SelectedDate     string
	SelectedRecord   *WorkDateRecord
	ProjectedEnd     string
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

type UserConfig struct {
	DefaultLunch     sql.NullInt16
	DefaultDayLength sql.NullString
	OffStart         sql.NullString
	OffEnd           sql.NullString
}

func openDBConnection() (conn *sql.DB, err error) {
	con, err := db.CreateConnection()
	if err != nil {
		return nil, err
	}

	return con, nil
}

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
			&workDateRecord.SickDay,
			&workDateRecord.DayLength,
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
		&workDateRecord.SickDay,
		&workDateRecord.DayLength,
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
		&userConfig.VacationStart,
		&userConfig.VacationEnd,
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

func WriteStart(selectedDate string, registeredTime string, dayLength string) error {
	con, err := openDBConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	query := "INSERT INTO timesheet(workdate, start_time, day_length) VALUES (?, ?, ?)"
	_, err = con.Exec(query, selectedDate, registeredTime, dayLength)
	if err != nil {
		return fmt.Errorf("failed to write start time to %s: %v", selectedDate, err)
	}

	return nil
}

func WriteNewBalance(selectedDate string, dayTotal string, dayBalance float64, totalBalance float64) error {
	con, err := openDBConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	query := "UPDATE timesheet " +
		"SET day_total = ?, day_balance = ?, moving_balance = ?" +
		"WHERE workdate = ?"

	_, err = con.Exec(query, dayTotal, dayBalance, totalBalance, selectedDate)
	if err != nil {
		return fmt.Errorf("failed to write new balance data: %v", err)
	}

	return nil
}

func UpdateStart(selectedDate string, registeredTime string) error {
	con, err := openDBConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	query := "UPDATE timesheet SET start_time = ? WHERE workdate = ?"
	_, err = con.Exec(query, registeredTime, selectedDate)
	if err != nil {
		return fmt.Errorf("failed to update start time%v", err)
	}

	return nil
}

func UpdateEnd(
	selectedDate string,
	registeredTime string,
	overtime bool,
	lunchDuration int16,
	additionalTime int16,
) error {
	con, err := openDBConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	query := "UPDATE timesheet " +
		"SET end_time = ?, overtime = ?, lunch_duration = ?, additional_time = ? " +
		"WHERE workdate = ?"
	_, err = con.Exec(query, registeredTime, overtime, lunchDuration, additionalTime, selectedDate)
	if err != nil {
		return fmt.Errorf("failed to update end time%v", err)
	}

	return nil
}

func UpdateLunch(selectedDate string, lunchDuration int16) error {
	con, err := openDBConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	query := "UPDATE timesheet SET lunch_duration = ? WHERE workdate = ?"
	_, err = con.Exec(query, lunchDuration, selectedDate)
	if err != nil {
		return fmt.Errorf("failed to update lunch duration%v", err)
	}

	return nil
}

func UpdateAdditionalTime(selectedDate string, additionalTime int16) error {
	con, err := openDBConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	query := "UPDATE timesheet SET additional_time = ? WHERE workdate = ?"
	_, err = con.Exec(query, additionalTime, selectedDate)
	if err != nil {
		return fmt.Errorf("failed to update additional time%v", err)
	}

	return nil
}

func UpdateDefaultLunch(lunchDuration int16) error {
	con, err := openDBConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	query := "UPDATE userconfig SET lunch_duration = ? WHERE ROWID = 1"
	_, err = con.Exec(query, lunchDuration)
	if err != nil {
		return fmt.Errorf("failed to update default lunch%v", err)
	}

	return nil
}

func UpdateDefaultLength(dayLength int16) error {
	con, err := openDBConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(con)

	query := "UPDATE userconfig SET length_of_day = ? WHERE ROWID = 1"
	_, err = con.Exec(query, dayLength)
	if err != nil {
		return fmt.Errorf("failed to update default lunch%v", err)
	}

	return nil
}
