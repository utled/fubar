package cli

import (
	"fTime/db"
	"fTime/home_helpers"
)

func Main() error {
	err := db.InitializeDB()
	if err != nil {
		return err
	}

	home_helpers.GetTimesheet()

	return nil
}
