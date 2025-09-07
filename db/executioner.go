package db

import (
	"database/sql"
	"fmt"
)

func GetMultiRecords(query string, startDate string, endDate string) (response *sql.Rows, err error) {
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

	response, err = db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return response, nil
}

func GetOneRecord(query string, constraint any) (response *sql.Row, err error) {
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

	response = db.QueryRow(query, constraint)

	return response, nil
}

func GetOneValue(query string, constraint any) (value string, err error) {
	db, err := createConnection()
	if err != nil {
		return "", err
	}
	defer func(db *sql.DB) {
		err := closeConnection(db)
		if err != nil {
			fmt.Println("failed to close connection:", err)
		}
	}(db)

	if constraint == nil {
		response := db.QueryRow(query)

		err = response.Scan(&value)
		if err != nil {
			return "", fmt.Errorf("failed to execute query: %v", err)
		}
	}

	return value, nil
}
