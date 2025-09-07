package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DefaultConfig struct {
	LunchDuration int
	LengthOfDay   string
}

func InitializeDB() error {
	db, err := createConnection()
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := closeConnection(db)
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
			workdate TEXT UNIQUE NOT NULL,
			start_time TEXT,
			end_time TEXT,
			lunch_duration INT,
			day_total TEXT,
			day_balance FLOAT,
			overtime INT,
			moving_balance FLOAT,
			additional_time INT,
			sick_day INT,
			day_length TEXT
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
	query := "SELECT lunch_duration, length_of_day FROM userconfig WHERE ROWID = ?;"
	queryResponse := db.QueryRow(query, 1)

	config := &DefaultConfig{}
	err := queryResponse.Scan(&config.LunchDuration, &config.LengthOfDay)
	if err != nil {
		if err == sql.ErrNoRows {
			config.LunchDuration = 40
			config.LengthOfDay = "08:00"
			_, err = db.Exec("INSERT INTO userconfig (lunch_duration, length_of_day) VALUES (?, ?)", config.LunchDuration, config.LengthOfDay)
			if err != nil {
				return fmt.Errorf("failed to write default user config: %v", err)
			}
		} else {
			return fmt.Errorf("failed to read user config: %v", err)
		}
	}

	return nil
}
