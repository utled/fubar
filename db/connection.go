package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

func CreateConnection() (db *sql.DB, err error) {
	dbUser := os.Getenv("FTIMEDBUser")
	dbPswd := os.Getenv("FTIMEDBPswd")
	dbHost := os.Getenv("FTIMEDBHost")

	cfg := mysql.NewConfig()
	cfg.User = dbUser
	cfg.Passwd = dbPswd
	cfg.Net = "tcp"
	cfg.Addr = dbHost
	cfg.DBName = "ftime"

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return db, err
	}

	return db, nil
}

func CloseConnection(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		return fmt.Errorf("faIled to close db connection: %v", err)
	}

	return nil
}
