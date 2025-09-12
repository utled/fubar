package actions

import (
	"fTime/helpers"
	"fTime/utils"
	"fmt"
	"time"
)

func RegisterWeekend(selectedDate string, userConfig *helpers.UserConfig) error {
	parsedDate, err := time.Parse(utils.DateLayout, selectedDate)
	if err != nil {
		return fmt.Errorf("failed to parse date: %v", err)
	}
	previousBalance, err := helpers.GetPreviousBalance(parsedDate)
	if err != nil {
		return fmt.Errorf("failed to get the previous balance: %v", err)
	}

	weekend := make([]string, 2)
	saturday := parsedDate.AddDate(0, 0, 1)
	weekend[0] = saturday.Format(utils.DateLayout)
	sunday := parsedDate.AddDate(0, 0, 2)
	weekend[1] = sunday.Format(utils.DateLayout)

	for _, date := range weekend {
		err = helpers.WriteWeekend(date, previousBalance, userConfig.DefaultDayLength.String)
		if err != nil {
			return fmt.Errorf("failed to write weekend: %v", err)
		}
	}

	return nil
}
