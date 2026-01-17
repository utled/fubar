package tui

import (
	"fmt"
	"fubar/data"
	"fubar/helpers"
	"fubar/registration"
	"fubar/utils"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type configMsg struct {
	config data.UserConfig
}

func (model *Model) fetchUserConfig() tea.Cmd {
	return func() tea.Msg {
		cfg, err := data.GetUserConfig()
		if err != nil {
			return errMsg(err)
		}
		return configMsg{config: cfg}
	}
}

type dateStateMsg struct {
	dateState data.ReportState
}

func (model *Model) calcProjectedEnd(selectedDateRecord *data.WorkDateRecord) string {
	startTime, err := helpers.ParseTimeObject(selectedDateRecord.StartTime.String)
	if err != nil {
		return ""
	}

	dayLength, err := helpers.ParseTimeObject(selectedDateRecord.DayLength.String)
	if err != nil {
		return ""
	}

	var lunchDuration time.Time
	if selectedDateRecord.LunchDuration.Valid {
		lunchDuration = dayLength.Add(time.Minute * time.Duration(selectedDateRecord.LunchDuration.Int16))
	} else {
		lunchDuration = dayLength.Add(time.Minute * time.Duration(model.userConfig.DefaultLunch.Int16))
	}

	addHour := time.Duration(lunchDuration.Hour()) * time.Hour
	addMinute := time.Duration(lunchDuration.Minute()) * time.Minute
	projectedEnd := startTime.Add(addHour + addMinute)

	return projectedEnd.Format(utils.TimeLayout)
}

func (model *Model) fetchDateState(selectedDate string) tea.Cmd {
	return func() tea.Msg {
		reportState := data.ReportState{}
		maxCompletedDate, maxDate, err := data.GetMaxDates()
		if err != nil {
			fmt.Println("failed to set new state.\n", err)
		}

		totalBalance, err := data.GetCurrentTotalBalance()
		if err != nil {
			fmt.Println("failed to set new state.\n", err)
		}

		previousCompleted, err := helpers.CheckPreviousCompletion(selectedDate, maxCompletedDate)
		if err != nil {
			fmt.Println("failed to set new state.\n", err)
		}

		recordExists, err := helpers.CheckIfDateExists(selectedDate, maxDate)
		if err != nil {
			fmt.Println("failed to set new state.\n", err)
		}

		var selectedDateRecord data.WorkDateRecord
		if recordExists {
			selectedDateRecord, err = data.GetOneWorkDateRecord(selectedDate)
			if err != nil {
				fmt.Println("failed to set new state.\n", err)
			}
		} else {
			selectedDateRecord = data.WorkDateRecord{
				WorkDate: selectedDate,
			}
		}

		reportState.SelectedDate = selectedDate
		reportState.ReportUpToDate = previousCompleted
		reportState.MaxDate = maxDate
		reportState.MaxCompletedDate = maxCompletedDate
		reportState.TotalBalance = totalBalance
		reportState.SelectedRecord = &selectedDateRecord
		projectedEnd := model.calcProjectedEnd(&selectedDateRecord)
		reportState.ProjectedEnd = projectedEnd
		return dateStateMsg{dateState: reportState}
	}
}

type tableDataMsg struct {
	Rows []*data.WorkDateRecord
}

func (model *Model) fetchTableData() tea.Cmd {
	return func() tea.Msg {
		minDate, err := data.GetMinDate()
		if err != nil {
			return errMsg(err)
		}
		allDateRecords, err := data.GetTimesheetRange(minDate, model.dateState.MaxDate)
		if err != nil {
			return errMsg(err)
		}
		return tableDataMsg{Rows: allDateRecords}
	}
}

func (model *Model) syncInputs() {
	userConfig := model.userConfig
	model.configInputFields[0].SetValue(strconv.Itoa(int(userConfig.DefaultLunch.Int16)))
	model.configInputFields[1].SetValue(userConfig.DefaultDayLength.String[:5])

	if userConfig.OffStart.String != "" {
		model.headerFields[3].Placeholder = "[y] Schedule*"
		model.scheduleInputFields[0].SetValue(userConfig.OffStart.String)
		model.scheduleInputFields[1].SetValue(userConfig.OffEnd.String)
		model.scheduleInputFields[2].SetValue(userConfig.OffType.String)
	} else {
		model.headerFields[3].Placeholder = "[y] Schedule"
		model.scheduleInputFields[0].SetValue("")
		model.scheduleInputFields[1].SetValue("")
		model.scheduleInputFields[2].SetValue("norm")
	}

	state := model.dateState

	if state.SelectedRecord.StartTime.Valid {
		model.dailyInputFields[idxStart].SetValue(state.SelectedRecord.StartTime.String[:5])
	} else {
		model.dailyInputFields[idxStart].SetValue("")
		model.dailyInputFields[idxStart].Placeholder = "07:00"
	}

	if state.SelectedRecord.LunchDuration.Valid {
		model.dailyInputFields[idxLunch].SetValue(strconv.Itoa(int(state.SelectedRecord.LunchDuration.Int16)))
	} else {
		model.dailyInputFields[idxLunch].SetValue("")
		model.dailyInputFields[idxLunch].Placeholder = strconv.Itoa(int(model.userConfig.DefaultLunch.Int16))
	}

	if state.SelectedRecord.EndTime.Valid {
		model.dailyInputFields[idxEnd].SetValue(state.SelectedRecord.EndTime.String[:5])
	} else {
		model.dailyInputFields[idxEnd].SetValue("")
		if state.ProjectedEnd == "" {
			model.dailyInputFields[idxEnd].Placeholder = "16:00"
		} else {
			model.dailyInputFields[idxEnd].Placeholder = state.ProjectedEnd[:5]
		}

	}

	if state.SelectedRecord.AdditionalTime.Valid {
		model.dailyInputFields[idxAdditional].SetValue(strconv.Itoa(int(state.SelectedRecord.AdditionalTime.Int16)))
	} else {
		model.dailyInputFields[idxAdditional].SetValue("")
		model.dailyInputFields[idxAdditional].Placeholder = "0"
	}

	if state.SelectedRecord.DayType.Valid {
		model.dailyInputFields[idxDayType].SetValue(state.SelectedRecord.DayType.String)
	} else {
		model.dailyInputFields[idxDayType].SetValue("norm")
	}

	if state.SelectedRecord.Overtime.Valid {
		model.dailyInputFields[idxOvertime].SetValue(strconv.FormatBool(state.SelectedRecord.Overtime.Bool))
	} else {
		model.dailyInputFields[idxOvertime].SetValue("false")
	}

	if state.SelectedRecord.DayTotal.Valid {
		model.dailyInputFields[idxDayTotal].SetValue(state.SelectedRecord.DayTotal.String[:5])
	} else {
		model.dailyInputFields[idxDayTotal].SetValue("")
		model.dailyInputFields[idxDayTotal].Placeholder = "00:00"
	}

	if state.SelectedRecord.DayBalance.Valid {
		model.dailyInputFields[idxDayBalance].SetValue(fmt.Sprintf("%.2f", state.SelectedRecord.DayBalance.Float64))
	} else {
		model.dailyInputFields[idxDayBalance].SetValue("")
		model.dailyInputFields[idxDayBalance].Placeholder = "0.00"
	}
}

type dayUpdateMsg struct{ updated bool }

func (model *Model) updateDay() tea.Cmd {
	return func() tea.Msg {
		dayUpdated := false

		existingState := model.dateState
		existingRecord := *existingState.SelectedRecord
		existingState.SelectedRecord = &existingRecord
		inputValues, err := helpers.CollectDailyInputs(
			model.dailyInputFields[idxStart].Value(),
			model.dailyInputFields[idxLunch].Value(),
			model.dailyInputFields[idxEnd].Value(),
			model.dailyInputFields[idxAdditional].Value(),
			model.dailyInputFields[idxDayType].Value(),
			model.dailyInputFields[idxOvertime].Value(),
		)
		if err != nil {
			return errMsg(err)
		}
		originalDayType := existingState.SelectedRecord.DayType.String
		originalOverTime := existingState.SelectedRecord.Overtime.Bool
		existingState.SelectedRecord.DayType.String = inputValues.DayType
		existingState.SelectedRecord.Overtime.Bool = inputValues.Overtime

		if inputValues.StartTime != existingRecord.StartTime.String {
			if inputValues.StartTime == "" {
				return errMsg(fmt.Errorf("can't set start time to 00:00.\nuse key Del to delete record"))
			}
			err = registration.RegisterStart(inputValues.StartTime, &existingState, &model.userConfig)
			if err != nil {
				return errMsg(fmt.Errorf("error registering start: %v", err))
			}
			dayUpdated = true
		}
		if inputValues.LunchDuration != int(existingRecord.LunchDuration.Int16) {
			err = registration.RegisterLunch(inputValues.LunchDuration, &existingState)
			if err != nil {
				return errMsg(err)
			}
			dayUpdated = true
		}
		if inputValues.EndTime != existingRecord.EndTime.String {
			err = registration.RegisterEnd(inputValues.EndTime, &existingState, &model.userConfig)
			if err != nil {
				return errMsg(err)
			}
			dayUpdated = true
		}
		if inputValues.AdditionalTime != int(existingRecord.AdditionalTime.Int16) {
			err = registration.RegisterAdditionalTime(inputValues.AdditionalTime, &existingState)
			if err != nil {
				return errMsg(err)
			}
			dayUpdated = true
		}
		existingState.SelectedRecord.DayType.String = originalDayType
		if inputValues.DayType != originalDayType {
			switch inputValues.DayType {
			case "norm":
				if !dayUpdated {
					dateExists, err := helpers.CheckIfDateExists(existingState.SelectedDate, existingState.MaxDate)
					if err != nil {
						return errMsg(err)
					}
					if dateExists {
						err = registration.RevertOffDay(&model.userConfig, &existingState)
						if err != nil {
							return errMsg(err)
						}
						dayUpdated = true
					}
				}
			case "off", "vac", "sic":
				err = registration.RegisterOffDay(&model.userConfig, &existingState, inputValues.DayType)
				if err != nil {
					return errMsg(err)
				}
				dayUpdated = true

			}
		}

		existingState.SelectedRecord.Overtime.Bool = originalOverTime
		if inputValues.Overtime != originalOverTime {
			if !dayUpdated {
				err := registration.RegisterOvertime(inputValues.Overtime, &existingState, &model.userConfig)
				if err != nil {
					return errMsg(err)
				}
				dayUpdated = true
			}
		}
		return dayUpdateMsg{updated: dayUpdated}
	}
}
