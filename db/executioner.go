package db

import (
	"database/sql"
	"fmt"
)

func GetTimesheetData() (response *sql.Rows, err error) {
	db, err := createConnection()
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		err := closeConnection(db)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(db)

	query := "SELECT * FROM timesheet;"
	response, err = db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return response, nil
}
