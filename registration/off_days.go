package registration

import (
	"fTime/data"
	"fTime/helpers"
	"fTime/utils"
	"fmt"
	"time"
)

func RegisterWeekend(saturday time.Time, userConfig *data.UserConfig, state *data.ReportState) error {
	weekend := make([]data.OffDay, 0)
	sunday := saturday.AddDate(0, 0, 1)
	weekend = append(weekend, data.OffDay{OffDate: saturday.Format(utils.DateLayout), OffType: "wknd"})
	weekend = append(weekend, data.OffDay{OffDate: sunday.Format(utils.DateLayout), OffType: "wknd"})
	previousBalance := state.SelectedRecord.MovingBalance.Float64

	err := data.WriteOffDays(&weekend, previousBalance, userConfig.DefaultDayLength.String)
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

func RegisterOffPeriod(nextDay time.Time, userConfig *data.UserConfig, state *data.ReportState) error {
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

	offPeriod := make([]data.OffDay, 0)

	for currentDate := nextDay; currentDate.Before(lastDay) || currentDate.Equal(lastDay); currentDate = currentDate.AddDate(0, 0, 1) {
		weekday := currentDate.Weekday()
		dateString := currentDate.Format(utils.DateLayout)
		if weekday == time.Saturday || weekday == time.Sunday {
			offPeriod = append(offPeriod, data.OffDay{OffDate: dateString, OffType: "wknd"})
		} else {
			offPeriod = append(offPeriod, data.OffDay{OffDate: dateString, OffType: userConfig.OffType.String})
		}
	}

	previousBalance := state.SelectedRecord.MovingBalance.Float64

	err = data.WriteOffDays(&offPeriod, previousBalance, userConfig.DefaultDayLength.String)
	if err != nil {
		return err
	}

	err = data.UpdateScheduledOff("", "", "")

	return nil
}

func registerFullOffDay(userConfig *data.UserConfig, state *data.ReportState, offType string) error {
	parsedDate, err := time.Parse(utils.DateLayout, state.SelectedDate)
	if err != nil {
		return fmt.Errorf("failed to parse selected date day: %v", err)
	}
	previousBalance, err := data.GetPreviousBalance(parsedDate)
	if err != nil {
		return err
	}
	offDay := []data.OffDay{{OffDate: state.SelectedDate, OffType: offType}}

	selectedBeforeMax, err := helpers.CheckDateBefore(state.SelectedDate, state.MaxCompletedDate)
	if err != nil {
		return err
	}

	if selectedBeforeMax {
		err = data.UpdateFullOffDay(&offDay, previousBalance, userConfig.DefaultDayLength.String)
		if err != nil {
			return err
		}
		state.SelectedRecord.MovingBalance.Float64 = previousBalance
		err = rebalanceSucceedingDates(state)
		if err != nil {
			return err
		}
	} else {
		err = data.WriteOffDays(&offDay, previousBalance, userConfig.DefaultDayLength.String)
		if err != nil {
			return err
		}
		state.SelectedRecord.MovingBalance.Float64 = previousBalance
		nextDay := parsedDate.AddDate(0, 0, 1)
		if userConfig.OffStart.String != "" {
			parsedScheduledStart, err := time.Parse(utils.DateLayout, userConfig.OffStart.String)
			if err != nil {
				return fmt.Errorf("failed to parse scheduled start date.%v", err)
			}
			if nextDay == parsedScheduledStart {
				err = RegisterOffPeriod(nextDay, userConfig, state)
				if err != nil {
					return err
				}
			}
		}
		if nextDay.Weekday() == time.Saturday {
			err = RegisterWeekend(nextDay, userConfig, state)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func registerPartialOffDay(state *data.ReportState, offType string) error {
	parsedDate, err := time.Parse(utils.DateLayout, state.SelectedDate)
	if err != nil {
		return fmt.Errorf("failed to parse selected date day: %v", err)
	}
	previousBalance, err := data.GetPreviousBalance(parsedDate)
	if err != nil {
		return err
	}
	offDay := []data.OffDay{{OffDate: state.SelectedDate, OffType: offType}}

	err = data.UpdatePartialOffDay(&offDay, previousBalance)
	if err != nil {
		return err
	}

	selectedBeforeMax, err := helpers.CheckDateBefore(state.SelectedDate, state.MaxCompletedDate)
	if err != nil {
		return err
	}

	if selectedBeforeMax {
		state.SelectedRecord.MovingBalance.Float64 = previousBalance
		err = rebalanceSucceedingDates(state)
		if err != nil {
			return err
		}
	}

	return nil
}

func RegisterOffDay(userConfig *data.UserConfig, state *data.ReportState, offType string) error {
	if !state.ReportUpToDate {
		return fmt.Errorf("can't register selected date.\nAll previous dates must be up to date.")
	}

	if state.SelectedRecord.EndTime.Valid {
		err := registerPartialOffDay(state, offType)
		if err != nil {
			return err
		}
		return nil
	}

	if state.SelectedRecord.StartTime.Valid {
		return fmt.Errorf("registered date must not be started, or must be ended.\n" +
			"use <delete> to remove date before or add an end time.")
	}

	err := registerFullOffDay(userConfig, state, offType)
	if err != nil {
		return err
	}

	return nil
}

func RevertOffDay(userConfig *data.UserConfig, state *data.ReportState) error {
	if state.SelectedRecord.DayType.String == "norm" {
		return fmt.Errorf("date is already flagged as a normal workday")
	}

	state.SelectedRecord.DayType.String = "norm"

	endTime := state.SelectedRecord.EndTime.String[:2] + state.SelectedRecord.EndTime.String[3:5]
	err := RegisterEnd(endTime, state.SelectedRecord.DayType.String, state, userConfig)
	if err != nil {
		return err
	}

	return nil
}
