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
			workdate VARCHAR(10) UNIQUE NOT NULL,
			day_type VARCHAR(4),
			start_time VARCHAR(8),
			end_time VARCHAR(8),
			lunch_duration INT,
			additional_time INT,
			overtime BOOLEAN,
			day_total VARCHAR(8),
			day_balance FLOAT,
			total_balance FLOAT,
			day_length VARCHAR(8)
		);`,
		/*		`CREATE TABLE IF NOT EXISTS userconfig (
				lunch_duration INT,
				length_of_day TEXT,
				vacation_start TEXT,
				vacation_end TEXT
			);`,*/
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
