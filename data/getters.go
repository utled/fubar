package data

import (
	"database/sql"
	"fmt"
	"fubar/db"
	"fubar/utils"
	"time"
)

func GetTimesheetRange(startDate string, endDate string) (timesheet []*WorkDateRecord, err error) {
	con, err := db.CreateConnection()
	if err != nil {
		return timesheet, err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := "SELECT * FROM timesheet WHERE workdate BETWEEN ? AND ?;"
	response, err := con.Query(query, startDate, endDate)
	if err != nil {
		return timesheet, fmt.Errorf("failed to execute query: %v", err)
	}

	for response.Next() {
		workDateRecord := &WorkDateRecord{}
		err := response.Scan(
			&workDateRecord.WorkDate,
			&workDateRecord.DayType,
			&workDateRecord.StartTime,
			&workDateRecord.EndTime,
			&workDateRecord.LunchDuration,
			&workDateRecord.AdditionalTime,
			&workDateRecord.Overtime,
			&workDateRecord.DayTotal,
			&workDateRecord.DayBalance,
			&workDateRecord.TotalBalance,
			&workDateRecord.DayLength,
		)
		if err != nil {
			return timesheet, fmt.Errorf("failed to serialize range of records to struct: %v", err)
		}
		timesheet = append(timesheet, workDateRecord)
	}

	return timesheet, nil
}

func GetOneWorkDateRecord(queryDate string) (record WorkDateRecord, err error) {
	con, err := db.CreateConnection()
	if err != nil {
		return WorkDateRecord{}, err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := "SELECT * FROM timesheet WHERE workdate = ?"
	response := con.QueryRow(query, queryDate)

	workDateRecord := &WorkDateRecord{}

	err = response.Scan(
		&workDateRecord.WorkDate,
		&workDateRecord.DayType,
		&workDateRecord.StartTime,
		&workDateRecord.EndTime,
		&workDateRecord.LunchDuration,
		&workDateRecord.AdditionalTime,
		&workDateRecord.Overtime,
		&workDateRecord.DayTotal,
		&workDateRecord.DayBalance,
		&workDateRecord.TotalBalance,
		&workDateRecord.DayLength,
	)
	if err != nil {
		return WorkDateRecord{}, fmt.Errorf("failed to serialize record to struct%v", err)
	}

	return *workDateRecord, nil

}

func GetMaxDates() (maxCompletedDate string, maxDate string, err error) {
	con, err := db.CreateConnection()
	if err != nil {
		return "", "", err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
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

func GetCurrentTotalBalance() (totalBalance float64, err error) {
	con, err := db.CreateConnection()
	if err != nil {
		return 0.0, err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := "SELECT total_balance FROM timesheet WHERE workdate = (SELECT MAX(workdate) FROM timesheet WHERE total_balance IS NOT NULL)"
	response := con.QueryRow(query)
	err = response.Scan(&totalBalance)
	if err != nil {
		return 0.0, fmt.Errorf("failed to read total balance: %v", err)
	}

	return totalBalance, nil
}

func GetUserConfig() (config UserConfig, err error) {
	con, err := db.CreateConnection()
	if err != nil {
		return UserConfig{}, err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := "SELECT * FROM userconfig WHERE id = 1;"
	response := con.QueryRow(query)

	userConfig := &UserConfig{}

	err = response.Scan(
		&userConfig.ID,
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
	con, err := db.CreateConnection()
	if err != nil {
		return 0.0, err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	previousDate := selectedDate.AddDate(0, 0, -1)
	previousDateString := previousDate.Format(utils.DateLayout)
	query := "SELECT total_balance FROM timesheet WHERE workdate = ?;"
	response := con.QueryRow(query, previousDateString)
	err = response.Scan(&previousBalance)
	if err != nil {
		return 0.0, fmt.Errorf("failed to read moving balance: %v", err)
	}

	return previousBalance, nil
}

func GetMonthlySummary(year int) (monthlySummary []*MonthStats, err error) {
	con, err := db.CreateConnection()
	if err != nil {
		return monthlySummary, err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := `SELECT
		MONTHNAME(workdate) as month,
		COUNT(CASE WHEN day_type NOT IN ('wknd', 'off') THEN 1 END) AS total_weekdays,
		COUNT(CASE WHEN day_type = 'norm' THEN 1 END) AS worked_days,
		SEC_TO_TIME(SUM(TIME_TO_SEC(day_total))) AS total_time_worked,
		COUNT(CASE WHEN day_type = 'vac' THEN 1 END) AS vacation_days,
		COUNT(CASE WHEN day_type = 'sic' THEN 1 END) AS sick_days,
		COUNT(CASE WHEN day_type = 'wknd' THEN 1 END) AS weekend_days,
		COUNT(CASE WHEN day_type = 'off' THEN 1 END) AS off_days,
		COUNT(CASE WHEN overtime = TRUE THEN 1 END) AS ot_days,
		SUM(CASE WHEN overtime = TRUE THEN day_balance END) AS total_ot
	FROM timesheet
	WHERE YEAR(workdate) = ?
	GROUP BY MONTHNAME(workdate);`

	response, err := con.Query(query, year)
	if err != nil {
		return monthlySummary, fmt.Errorf("failed to execute query: %v", err)
	}

	for response.Next() {
		monthStats := &MonthStats{}
		err := response.Scan(
			&monthStats.Month,
			&monthStats.TotalWeekDays,
			&monthStats.WorkedDays,
			&monthStats.WorkedTime,
			&monthStats.VacationDays,
			&monthStats.SickDays,
			&monthStats.WeekendDays,
			&monthStats.OffDays,
			&monthStats.OverTimeDays,
			&monthStats.TotalOvertime,
		)
		if err != nil {
			return monthlySummary, fmt.Errorf("failed to serialize range of records to struct: %v", err)
		}
		monthlySummary = append(monthlySummary, monthStats)
	}

	return monthlySummary, nil
}

func GetFullStatistics(startDate string, endDate string) (fullStatistics *FullStats, err error) {
	con, err := db.CreateConnection()
	if err != nil {
		return fullStatistics, err
	}
	defer func(con *sql.DB) {
		err = db.CloseConnection(con)
		if err != nil {
			fmt.Println(err)
		}
	}(con)

	query := `select
		COUNT(CASE WHEN day_type = 'norm' THEN 1 END) AS worked_days,
		COUNT(CASE WHEN day_type not in ('wknd', 'off') THEN 1 END) AS total_weekdays,
		SEC_TO_TIME(sum(TIME_TO_SEC(day_total))) AS total_time_worked,
		SEC_TO_TIME(
			 round(
					 avg(CASE WHEN day_type = 'norm' THEN TIME_TO_SEC(start_time) END)
			 )
		) AS avg_start,
		SEC_TO_TIME(
			 ROUND(
					 AVG(CASE WHEN day_type = 'norm' THEN TIME_TO_SEC(end_time) END)
			 )
		) AS avg_end,
		AVG(CASE WHEN day_type = 'norm' THEN lunch_duration END) AS avg_lunch,
		COUNT(CASE WHEN day_type = 'sic' THEN 1 END) AS sick_days,
		COUNT(CASE WHEN day_type = 'vac' THEN 1 END) AS vacation_days,
		COUNT(CASE WHEN overtime = TRUE THEN 1 END) AS ot_days,
		SUM(CASE WHEN overtime = TRUE THEN day_balance END) AS total_ot,
		AVG(CASE WHEN overtime = TRUE THEN day_balance END) AS avg_ot
	FROM timesheet
	WHERE workdate BETWEEN ? AND ?;`

	response := con.QueryRow(query, startDate, endDate)

	fullStatistics = &FullStats{}

	err = response.Scan(
		&fullStatistics.WorkedDays,
		&fullStatistics.TotalWeekDays,
		&fullStatistics.WorkedTime,
		&fullStatistics.AvgStart,
		&fullStatistics.AvgEnd,
		&fullStatistics.AvgLunch,
		&fullStatistics.SickDays,
		&fullStatistics.VacationDays,
		&fullStatistics.OverTimeDays,
		&fullStatistics.TotalOvertime,
		&fullStatistics.AvgOvertime,
	)
	if err != nil {
		return fullStatistics, fmt.Errorf("failed to serialize record to struct%v", err)
	}
	return fullStatistics, nil
}
