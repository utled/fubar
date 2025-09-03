package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func createConnection() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", "/home/utled/GolandProjects/fTime/db/fTime.db")
	if err != nil {
		return db, fmt.Errorf("failed to connect to db: %v", err)
	}

	return db, nil
}

func closeConnection(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		return fmt.Errorf("faIled to close db connection: %v", err)
	}

	return nil
}
