package cli

import (
	"fTime/db"
	"fTime/home_helpers"
	"fmt"
)

func Main() error {

	err := db.InitializeDB()
	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}

	home_helpers.InitClearFunctions()

	home()

	return nil
}
