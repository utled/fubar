package actions

import (
	"fTime/helpers"
	"fTime/utils"
	"fmt"
	"time"
)

func RegisterWeekend(saturday time.Time, userConfig *helpers.UserConfig, state *helpers.ReportState) error {
	weekend := make([]helpers.OffDay, 0)
	sunday := saturday.AddDate(0, 0, 1)
	weekend = append(weekend, helpers.OffDay{OffDate: saturday.Format(utils.DateLayout), OffType: "wknd"})
	weekend = append(weekend, helpers.OffDay{OffDate: sunday.Format(utils.DateLayout), OffType: "wknd"})
	previousBalance := state.SelectedRecord.MovingBalance.Float64

	err := helpers.WriteOffDays(&weekend, previousBalance, userConfig.DefaultDayLength.String)
	if err != nil {
		return err
	}

	if userConfig.OffStart.String != "" {
		parsedStart, err := time.Parse(utils.DateLayout, userConfig.OffStart.String)
		if err != nil {
			return fmt.Errorf("failed to parse start date: %v", err)
		}
		monday := sunday.AddDate(0, 0, 1)
		if monday.Equal(parsedStart) {
			err = RegisterOffPeriod(monday, userConfig, state)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func RegisterOffPeriod(nextDay time.Time, userConfig *helpers.UserConfig, state *helpers.ReportState) error {
	if userConfig.OffStart.String != "" {
		parsedStart, err := time.Parse(utils.DateLayout, userConfig.OffStart.String)
		if err != nil {
			return fmt.Errorf("failed to parse scheduled start: %v", err)
		}
		if parsedStart == nextDay.AddDate(0, 0, 1) {
		}
	}

	lastDay, err := time.Parse(utils.DateLayout, userConfig.OffEnd.String)
	if err != nil {
		return fmt.Errorf("failed to parse last scheduled day: %v", err)
	}

	offPeriod := make([]helpers.OffDay, 0)

	for currentDate := nextDay; currentDate.Before(lastDay) || currentDate.Equal(lastDay); currentDate = currentDate.AddDate(0, 0, 1) {
		weekday := currentDate.Weekday()
		dateString := currentDate.Format(utils.DateLayout)
		if weekday == time.Saturday || weekday == time.Sunday {
			offPeriod = append(offPeriod, helpers.OffDay{OffDate: dateString, OffType: "wknd"})
		} else {
			offPeriod = append(offPeriod, helpers.OffDay{OffDate: dateString, OffType: userConfig.OffType.String})
		}
	}

	previousBalance := state.SelectedRecord.MovingBalance.Float64

	err = helpers.WriteOffDays(&offPeriod, previousBalance, userConfig.DefaultDayLength.String)
	if err != nil {
		return err
	}

	err = helpers.UpdateScheduledOff("", "", "")

	return nil
}
