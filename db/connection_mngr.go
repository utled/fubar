package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"os"
)

type credentials struct {
	User string `json:"User"`
	Pswd string `json:"Pswd"`
	Host string `json:"Host"`
}

func CreateConnection() (db *sql.DB, err error) {
	var dbCredentials credentials
	credentialsJson := os.Getenv("FTIMEDB")
	err = json.Unmarshal([]byte(credentialsJson), &dbCredentials)
	if err != nil {
		return db, err
	}

	cfg := mysql.NewConfig()
	cfg.User = dbCredentials.User
	cfg.Passwd = dbCredentials.Pswd
	cfg.Net = "tcp"
	cfg.Addr = dbCredentials.Host
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
