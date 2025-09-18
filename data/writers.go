package data

import (
	"database/sql"
	"fTime/db"
	"fmt"
)

func WriteStart(selectedDate string, registeredTime string, dayLength string) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
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
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := "UPDATE timesheet " +
		"SET day_total = ?, day_balance = ?, total_balance = ?" +
		"WHERE workdate = ?"

	_, err = con.Exec(query, dayTotal, dayBalance, totalBalance, selectedDate)
	if err != nil {
		return fmt.Errorf("failed to write new balance data: %v", err)
	}

	return nil
}

func WriteOffDays(offPeriod *[]OffDay, totalBalance float64, defaultDayLength string) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := `INSERT INTO timesheet(
                      workdate, 
                      day_type,
                      start_time,
                      end_time,
                      lunch_duration,
                      additional_time,
                      overtime,
                      day_total,
                      day_balance,
                      total_balance,
                      day_length)
                      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, day := range *offPeriod {
		_, err = con.Exec(
			query,
			day.OffDate,
			day.OffType,
			"00:00:00",
			"00:00:00",
			0,
			0,
			false,
			"00:00:00",
			0.0,
			totalBalance,
			defaultDayLength)
		if err != nil {
			return fmt.Errorf("failed to write weekend: %v", err)
		}
	}

	return nil
}

func WriteBackflush(backflushRange *[]WorkDateRecord) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := `INSERT INTO timesheet(
                      workdate,
                      day_type,
                      start_time,
                      end_time,
                      lunch_duration,
                      additional_time,
                      overtime,
                      day_total,
                      day_balance,
                      total_balance,
                      day_length)
                      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	for _, day := range *backflushRange {
		_, err = con.Exec(
			query,
			day.WorkDate,
			day.DayType,
			day.StartTime,
			day.EndTime,
			day.LunchDuration,
			day.AdditionalTime,
			day.Overtime,
			day.DayTotal,
			day.DayBalance,
			day.TotalBalance,
			day.DayLength)
		if err != nil {
			return fmt.Errorf("failed to write backflush records: %v", err)
		}
	}

	return nil
}

func UpdateStart(selectedDate string, registeredTime string) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
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
	dayType string,
) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := `UPDATE timesheet 
SET end_time = ?, overtime = ?, lunch_duration = ?, additional_time = ?, day_type = ? 
WHERE workdate = ?`
	_, err = con.Exec(query, registeredTime, overtime, lunchDuration, additionalTime, dayType, selectedDate)
	if err != nil {
		return fmt.Errorf("failed to update end time%v", err)
	}

	return nil
}

func UpdateLunch(selectedDate string, lunchDuration int16) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
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
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
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
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := "UPDATE userconfig SET default_lunch = ? WHERE id = 1"
	_, err = con.Exec(query, lunchDuration)
	if err != nil {
		return fmt.Errorf("failed to update default lunch%v", err)
	}

	return nil
}

func UpdateDefaultLength(dayLength int16) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := "UPDATE userconfig SET default_day_length = ? WHERE id = 1"
	_, err = con.Exec(query, dayLength)
	if err != nil {
		return fmt.Errorf("failed to update default lunch%v", err)
	}

	return nil
}

func UpdateScheduledOff(offStart string, offEnd string, offType string) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := `UPDATE userconfig 
SET scheduled_off_start = ?, scheduled_off_end = ?, scheduled_off_type = ? 
WHERE ROWID = 1`
	_, err = con.Exec(query, offStart, offEnd, offType)
	if err != nil {
		return fmt.Errorf("failed to update scheduled off period: %v", err)
	}

	return nil
}

func UpdateFullOffDay(offPeriod *[]OffDay, totalBalance float64, defaultDayLength string) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := `UPDATE timesheet
    SET day_type = ?,
    	start_time = ?, 
        end_time = ?, 
        lunch_duration = ?, 
        additional_time = ?, 
        overtime = ?,
        day_total = ?, 
        day_balance = ?,  
        total_balance = ?, 
        day_length = ?
    WHERE workdate = ?`

	for _, day := range *offPeriod {
		_, err = con.Exec(
			query,
			day.OffType,
			"00:00:00",
			"00:00:00",
			0,
			0,
			false,
			"00:00:00",
			0.0,
			totalBalance,
			defaultDayLength,
			day.OffDate)
		if err != nil {
			return fmt.Errorf("failed to write weekend: %v", err)
		}
	}

	return nil
}

func UpdatePartialOffDay(offPeriod *[]OffDay, totalBalance float64) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := `UPDATE timesheet
    SET total_balance = ?,
        day_type = ?
    WHERE workdate = ?`

	for _, day := range *offPeriod {
		_, err = con.Exec(
			query,
			totalBalance,
			day.OffType,
			day.OffDate)
		if err != nil {
			return fmt.Errorf("failed to write weekend: %v", err)
		}
	}

	return nil
}

func UpdateTotalBalance(dateRange *[]string, previousBalance float64) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := `UPDATE timesheet SET total_balance = day_balance + ? WHERE workdate = ?`
	for _, day := range *dateRange {
		_, err = con.Exec(query, previousBalance, day)
		if err != nil {
			return fmt.Errorf("failed to update balance: %v", err)
		}
	}

	return nil
}

func DeleteRecord(selectedDate string) error {
	con, err := db.CreateConnection()
	if err != nil {
		return err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := "DELETE FROM timesheet WHERE workdate = ?"
	_, err = con.Exec(query, selectedDate)
	if err != nil {
		return fmt.Errorf("failed to delete record: %v", err)
	}

	return nil
}
