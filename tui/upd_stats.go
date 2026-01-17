package tui

import (
	"fmt"
	"fubar/data"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) updateStats(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statsYearMsg:
		model.statsDetails.minYear = msg.minYear
		model.statsDetails.maxYear = msg.maxYear
		maxYearString := strconv.Itoa(model.statsDetails.maxYear)
		model.statsDetails.yearSelection.SetValue(maxYearString)

		var cmds []tea.Cmd
		var cmd tea.Cmd
		model.statsDetails.yearSelection, cmd = model.statsDetails.yearSelection.Update(msg)
		cmds = append(cmds, cmd)
		cmds = append(cmds, model.generateStatsGraph())
		cmds = append(cmds, model.fetchAllSumData())
		cmds = append(cmds, model.fetchTableData())
		cmds = append(cmds, model.fetchAllSumData())
		return model, tea.Batch(cmds...)
	case statsGraphMsg:
		model.statsDetails.graphArea.SetContent(msg.graphString)
		model.statsDetails.graphArea.GotoTop()
		return model, nil
	case statsTableDataMsg:
		var rows []table.Row
		for _, row := range msg.tableData {
			tableRow := generateStatsTableRow(row, &model.statsDetails.tableTotals)
			rows = append(rows, tableRow)
		}
		missingMonths := 12 - len(rows)
		for i := 12 - missingMonths; i <= missingMonths; i++ {
			blankRow := generateBlankRow()
			currentMonth := time.Month(i + 1).String()
			blankRow[0] = currentMonth
			rows = append(rows, blankRow)
		}
		totalsRow := generateBottomRows(&model.statsDetails.tableTotals)
		for _, row := range totalsRow {
			rows = append(rows, row)
		}

		model.statsDetails.table.SetRows(rows)
		model.statsDetails.table.GotoTop()
		var cmds []tea.Cmd
		var cmd tea.Cmd
		model.statsDetails.table, cmd = model.statsDetails.table.Update(msg)
		cmds = append(cmds, cmd)
		return model, tea.Batch(cmds...)
	case statsAllSumDataMsg:
		model.statsDetails.allSumFields[idxWorkedDays].SetValue(strconv.Itoa(msg.fieldData.WorkedDays))
		model.statsDetails.allSumFields[idxWeekdays].SetValue(strconv.Itoa(msg.fieldData.TotalWeekDays))
		model.statsDetails.allSumFields[idxWorkedTime].SetValue(msg.fieldData.WorkedTime)
		model.statsDetails.allSumFields[idxAvgStart].SetValue(msg.fieldData.AvgStart)
		model.statsDetails.allSumFields[idxAvgEnd].SetValue(msg.fieldData.AvgEnd)
		model.statsDetails.allSumFields[idxAvgLunch].SetValue(fmt.Sprintf("%.2f", msg.fieldData.AvgLunch))
		model.statsDetails.allSumFields[idxSickDays].SetValue(strconv.Itoa(msg.fieldData.SickDays))
		model.statsDetails.allSumFields[idxVacDays].SetValue(strconv.Itoa(msg.fieldData.VacationDays))
		model.statsDetails.allSumFields[idxOTDays].SetValue(strconv.Itoa(msg.fieldData.OverTimeDays))
		model.statsDetails.allSumFields[idxTotalOT].SetValue(fmt.Sprintf("%.2f", msg.fieldData.TotalOvertime.Float64))
		model.statsDetails.allSumFields[idxAvgOT].SetValue(fmt.Sprintf("%.2f", msg.fieldData.AvgOvertime.Float64))
		return model, nil
	case statsMonthSumDataMsg:
		model.statsDetails.monthSumFields[0].SetValue(msg.fieldData.AvgStart)
		model.statsDetails.monthSumFields[1].SetValue(msg.fieldData.AvgEnd)
		if msg.fieldData.AvgLunch > 0 {
			model.statsDetails.monthSumFields[2].SetValue(fmt.Sprintf("%.2f", msg.fieldData.AvgLunch))
		} else {
			model.statsDetails.monthSumFields[2].SetValue("")
		}
		if msg.fieldData.AvgOvertime.Float64 > 0 {
			model.statsDetails.monthSumFields[3].SetValue(fmt.Sprintf("%.2f", msg.fieldData.AvgOvertime.Float64))
		} else {
			model.statsDetails.monthSumFields[3].SetValue("")
		}

		return model, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		case "d":
			model.state = stateDaily
			return model, nil
		case "v":
			model.statsDetails.displayType = (model.statsDetails.displayType + 1) % 2
			var cmds []tea.Cmd
			switch model.statsDetails.displayType {
			case graphDisplay:
				cmds = append(cmds, model.generateStatsGraph())
				cmds = append(cmds, model.fetchAllSumData())
			case tableDisplay:
				cmds = append(cmds, model.fetchStatsTableData())
				cmds = append(cmds, model.fetchMonthSumData("January"))
			}
			return model, tea.Batch(cmds...)
		case "j", "k":
			var cmds []tea.Cmd
			var cmd tea.Cmd
			switch msg.String() {
			case "j":
				model.statsDetails.table, cmd = model.statsDetails.table.Update(tea.KeyMsg{Type: tea.KeyUp})
				cmds = append(cmds, cmd)
			case "k":
				model.statsDetails.table, cmd = model.statsDetails.table.Update(tea.KeyMsg{Type: tea.KeyDown})
				cmds = append(cmds, cmd)
			}
			cmds = append(cmds, model.fetchMonthSumData(model.statsDetails.table.SelectedRow()[0]))
			return model, tea.Batch(cmds...)
		case "n", "h", "p", "l":
			selectedYear, _ := strconv.Atoi(model.statsDetails.yearSelection.Value())
			switch msg.String() {
			case "n", "l":
				if selectedYear == model.statsDetails.maxYear {
					model.statsDetails.yearSelection.SetValue(strconv.Itoa(model.statsDetails.minYear))
				} else {
					selectedYear++
					model.statsDetails.yearSelection.SetValue(strconv.Itoa(selectedYear))
				}
			case "p", "h":
				if selectedYear == model.statsDetails.minYear {
					model.statsDetails.yearSelection.SetValue(strconv.Itoa(model.statsDetails.maxYear))
				} else {
					selectedYear--
					model.statsDetails.yearSelection.SetValue(strconv.Itoa(selectedYear))
				}
			}
			model.statsDetails.tableTotals = data.MonthStats{}
			var cmds []tea.Cmd
			switch model.statsDetails.displayType {
			case graphDisplay:
				cmds = append(cmds, model.generateStatsGraph())
			case tableDisplay:
				cmds = append(cmds, model.fetchStatsTableData())
				cmds = append(cmds, model.fetchMonthSumData("January"))
			}
			return model, tea.Batch(cmds...)
		case "b":
			model.prevState = model.state
			model.backflushInputFields.SetValue("norm")
			model.state = stateBackflush
			return model, nil
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
		}
		switch msg.Type {
		case tea.KeyUp, tea.KeyDown:
			var cmds []tea.Cmd
			var cmd tea.Cmd
			model.statsDetails.table, cmd = model.statsDetails.table.Update(msg)
			cmds = append(cmds, cmd)
			cmds = append(cmds, model.fetchMonthSumData(model.statsDetails.table.SelectedRow()[0]))
			return model, tea.Batch(cmds...)
		case tea.KeyLeft, tea.KeyRight, tea.KeyTab:
			var cmds []tea.Cmd
			selectedYear, _ := strconv.Atoi(model.statsDetails.yearSelection.Value())
			switch msg.Type {
			case tea.KeyLeft:
				if selectedYear == model.statsDetails.minYear {
					model.statsDetails.yearSelection.SetValue(strconv.Itoa(model.statsDetails.maxYear))
				} else {
					selectedYear--
					model.statsDetails.yearSelection.SetValue(strconv.Itoa(selectedYear))
				}
			default:
				if selectedYear == model.statsDetails.maxYear {
					model.statsDetails.yearSelection.SetValue(strconv.Itoa(model.statsDetails.minYear))
				} else {
					selectedYear++
					model.statsDetails.yearSelection.SetValue(strconv.Itoa(selectedYear))
				}
			}
			model.statsDetails.tableTotals = data.MonthStats{}
			switch model.statsDetails.displayType {
			case graphDisplay:
				cmds = append(cmds, model.generateStatsGraph())
			case tableDisplay:
				cmds = append(cmds, model.fetchStatsTableData())
				cmds = append(cmds, model.fetchMonthSumData("January"))
			}
			return model, tea.Batch(cmds...)
		case tea.KeyEnter:
		case tea.KeyEsc:
			model.state = stateDaily
			return model, nil
		}
	}
	var cmd tea.Cmd
	model.statsDetails.yearSelection, cmd = model.statsDetails.yearSelection.Update(msg)
	return model, cmd
}
