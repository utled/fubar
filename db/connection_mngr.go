package db

import (
	"database/sql"
	"fmt"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

func createConnection() (db *sql.DB, err error) {
	var dbPath string
	if runtime.GOOS == "windows" {
		dbPath = `C:\Users\utled\GolandProjects\fTime\db\fTime.db`
	} else {
		dbPath = "/home/utled/GolandProjects/fTime/db/fTime.db"
	}
	db, err = sql.Open("sqlite3", dbPath)
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
