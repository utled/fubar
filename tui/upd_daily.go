package tui

import (
	"fmt"
	"fubar/helpers"
	"fubar/utils"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) updateDaily(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case configMsg:
		model.userConfig = msg.config
		return model, model.fetchDateState(model.dateState.SelectedDate)
	case dateStateMsg:
		var cmds []tea.Cmd
		model.dateState = msg.dateState
		model.syncInputs()
		var cmd tea.Cmd
		model.dailyInputFields[model.dailyInputFocus], cmd = model.dailyInputFields[model.dailyInputFocus].Update(msg)
		cmds = append(cmds, cmd)
		cmds = append(cmds, model.fetchTableData())
		return model, tea.Batch(cmds...)
	case tableDataMsg:
		model.timesheet = msg.Rows
		var rows []table.Row
		var selectedIdx int
		for idx, row := range msg.Rows {
			tableRow := generateTableRow(row)
			rows = append(rows, tableRow)
			if tableRow[0] == model.dateState.SelectedDate {
				selectedIdx = idx
			}
		}
		model.dateTable.SetRows(rows)
		model.dateTable.SetCursor(selectedIdx)
		model.dateTable.Update(msg)
		return model, nil
	case dayUpdateMsg:
		model.dailyInputFields[model.dailyInputFocus].Blur()
		if !msg.updated {
			model.statusLabel = "no changes made to selected day"
			return model, clearStatusTimer()
		}
		model.statusLabel = "day updated successfully"
		var cmds []tea.Cmd
		cmds = append(cmds, model.fetchDateState(model.dateState.SelectedDate))
		cmds = append(cmds, clearStatusTimer())
		return model, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		case "ctrl+x":
			model.dailyInputFields[model.dailyInputFocus].SetValue("")
		case "b":
			model.prevState = model.state
			model.backflushInputFields.SetValue("norm")
			model.state = stateBackflush
			return model, nil
		case "x":
			model.statsDetails.displayType = graphDisplay
			model.state = stateStatistics
			return model, model.collectYearRange()
		case "y":
			model.prevState = model.state
			model.otherInputFocus = 0
			model.scheduleInputFields[model.otherInputFocus].Focus()
			model.scheduleInputFields[model.otherInputFocus].CursorEnd()
			model.state = stateSchedule
			return model, nil
		case "z":
			model.prevState = model.state
			model.otherInputFocus = 0
			model.configInputFields[model.otherInputFocus].Focus()
			model.configInputFields[model.otherInputFocus].CursorEnd()
			model.state = stateConfig
			return model, nil

		case "s", "l", "e", "a":
			var cmds []tea.Cmd
			model.dailyInputFields[model.dailyInputFocus].Blur()
			switch msg.String() {
			case "s":
				model.dailyInputFocus = idxStart
			case "l":
				model.dailyInputFocus = idxLunch
			case "e":
				model.dailyInputFocus = idxEnd
			case "a":
				model.dailyInputFocus = idxAdditional
			}
			cmds = append(cmds, model.dailyInputFields[model.dailyInputFocus].Focus())
			model.dailyInputFields[model.dailyInputFocus].CursorEnd()
			return model, tea.Batch(cmds...)

		case "t", "o":
			switch msg.String() {
			case "t":
				nextDayType := ParseDayType(model.dailyInputFields[idxDayType].Value()).Next().String()
				model.dailyInputFields[idxDayType].SetValue(nextDayType)
				return model, nil
			case "o":
				if model.dailyInputFields[idxOvertime].Value() == "false" {
					model.dailyInputFields[idxOvertime].SetValue("true")
				} else {
					model.dailyInputFields[idxOvertime].SetValue("false")
				}
				return model, nil
			}

		case "c", "n", "p":
			switch msg.String() {
			case "c":
				model.dateState.SelectedDate = time.Now().Format(utils.DateLayout)
				model.dateTable.GotoTop()
				return model, model.fetchDateState(model.dateState.SelectedDate)
			case "n":
				nextDate, err := helpers.NextDateFromString(model.dateState.SelectedDate)
				if err != nil {
					fmt.Println(err)
				}
				model.dateState.SelectedDate = nextDate
				return model, model.fetchDateState(model.dateState.SelectedDate)
			case "p":
				previousDate, err := helpers.PreviousDateFromString(model.dateState.SelectedDate)
				if err != nil {
					fmt.Println(err)
				}
				model.dateState.SelectedDate = previousDate
				return model, model.fetchDateState(model.dateState.SelectedDate)
			}
		}

		switch msg.Type {
		case tea.KeyUp, tea.KeyDown:
			var cmds []tea.Cmd
			var tableCmd tea.Cmd
			model.dateTable, tableCmd = model.dateTable.Update(msg)
			cmds = append(cmds, tableCmd)
			model.dateState.SelectedDate = model.dateTable.SelectedRow()[0]
			cmds = append(cmds, model.fetchDateState(model.dateState.SelectedDate))
			return model, tea.Batch(cmds...)
		case tea.KeyLeft, tea.KeyRight, tea.KeyTab:
			var cmds []tea.Cmd
			model.dailyInputFields[model.dailyInputFocus].Blur()
			focusCycleLimit := 4
			if msg.Type == tea.KeyLeft {
				model.dailyInputFocus = (model.dailyInputFocus - 1 + focusCycleLimit) % focusCycleLimit
			} else {
				model.dailyInputFocus = (model.dailyInputFocus + 1) % focusCycleLimit
			}
			cmds = append(cmds, model.dailyInputFields[model.dailyInputFocus].Focus())
			return model, tea.Batch(cmds...)
		case tea.KeyEnter:
			return model, model.updateDay()
		case tea.KeyEsc:
			var cmd tea.Cmd
			model.syncInputs()
			model.dailyInputFields[model.dailyInputFocus], cmd = model.dailyInputFields[model.dailyInputFocus].Update(msg)
			return model, cmd
		case tea.KeyDelete:
			model.confirmationDetails.confirmationType = deleteDate
			model.confirmationDetails.confirmationMsg = fmt.Sprint(
				"Please confirm deletion of date:\n",
				model.dateState.SelectedDate,
			)
			model.state = stateConfirm
			return model, nil
		}
	}

	var cmd tea.Cmd
	model.dailyInputFields[model.dailyInputFocus], cmd = model.dailyInputFields[model.dailyInputFocus].Update(msg)
	return model, cmd
}
