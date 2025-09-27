package db

import (
	"database/sql"
	"fmt"
)

type DefaultConfig struct {
	DefaultLunch     int
	DefaultDayLength string
}

func InitializeDB() error {
	db, err := CreateConnection()
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := CloseConnection(db)
		if err != nil {
			fmt.Println(err)
		}
	}(db)

	err = createTables(db)
	if err != nil {
		return err
	}

	err = writeDefaultConfig(db)
	if err != nil {
		return err
	}

	return nil
}

func createTables(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS timesheet (
			workdate DATE PRIMARY KEY,
			day_type VARCHAR(4),
			start_time TIME,
			end_time TIME,
			lunch_duration INT,
			additional_time INT,
			overtime BOOLEAN,
			day_total TIME,
			day_balance FLOAT,
			total_balance FLOAT,
			day_length TIME
		);`,
		`CREATE TABLE IF NOT EXISTS userconfig (
    			ID INT AUTO_INCREMENT PRIMARY KEY,
				default_lunch INT,
				default_day_length TIME,
				scheduled_off_start VARCHAR(10),
				scheduled_off_end VARCHAR(10),
				scheduled_off_type VARCHAR(3)
			);`,
	}

	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			return fmt.Errorf("failed to create table %s \n:%v", statement, err)
		}
	}

	return nil
}

func writeDefaultConfig(db *sql.DB) error {
	query := "SELECT default_lunch, default_day_length FROM userconfig WHERE id = 1;"
	queryResponse := db.QueryRow(query)

	config := &DefaultConfig{}
	err := queryResponse.Scan(&config.DefaultLunch, &config.DefaultDayLength)
	if err != nil {
		if err == sql.ErrNoRows {
			config.DefaultLunch = 40
			config.DefaultDayLength = "08:00"
			_, err = db.Exec("INSERT INTO userconfig (default_lunch, default_day_length) VALUES (?, ?)", config.DefaultLunch, config.DefaultDayLength)
			if err != nil {
				return fmt.Errorf("failed to write default user config: %v", err)
			}
		} else {
			return fmt.Errorf("failed to read user config: %v", err)
		}
	}

	return nil
}
