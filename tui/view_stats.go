package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func (model *Model) statsHeaderView() string {
	model.headerFields[0].SetValue("")
	model.headerFields[0].Placeholder = "  [d] Daily"
	model.headerFields[2].SetValue("    Stats")
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		model.inputStyle.InputField.Render(model.headerFields[0].View()),
		lipgloss.NewStyle().BorderForeground(lipgloss.Color("238")).BorderStyle(lipgloss.NormalBorder()).PaddingLeft(0).Render(model.headerFields[1].View()),
		model.inputStyle.InputField.Render(model.headerFields[2].View()),
		model.inputStyle.InputField.Render(model.headerFields[3].View()),
		model.inputStyle.InputField.Render(model.headerFields[4].View()),
	)
}

func (model *Model) statsYearSelectionView() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		"<-[p] ",
		model.inputStyle.InputField.PaddingRight(1).Render(model.statsDetails.yearSelection.View()),
		" [n]->   ",
	)
}

func (model *Model) statsSumView() string {
	var fieldData []textinput.Model
	switch model.statsDetails.displayType {
	case graphDisplay:
		fieldData = model.statsDetails.yearSumFields
	case tableDisplay:
		fieldData = model.statsDetails.allSumFields
	}
	sumOne := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%25s", "Worked Days: "), fieldData[idxWorkedDays].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%25s", "Weekdays: "), fieldData[idxWeekdays].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%25s", "Worked Time: "), fieldData[idxWorkedTime].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%25s", "Total OT: "), fieldData[idxTotalOT].View()),
	)
	sumTwo := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Avg Start: "), fieldData[idxAvgStart].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Avg End: "), fieldData[idxAvgEnd].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Avg Lunch: "), fieldData[idxAvgLunch].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Avg OT: "), fieldData[idxAvgOT].View()),
	)
	sumThree := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Sic Days: "), fieldData[idxSickDays].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Vac Days: "), fieldData[idxVacDays].View()),
		"",
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "OT Days: "), fieldData[idxOTDays].View()),
	)
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		sumOne,
		sumTwo,
		sumThree,
		"      ",
	)
}

func (model *Model) statsGraphView() string {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			model.statsDetails.graphArea.View(),
			"\n",
			lipgloss.NewStyle().Bold(true).Underline(true).Render("Selected Year"),
			model.statsSumView(),
		)
}

func (model *Model) statsMonthAvgView() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%-12s", "Avg Start")), model.statsDetails.monthAvgFields[0].View()),
		lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%-12s", "Avg End")), model.statsDetails.monthAvgFields[1].View()),
		lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%-12s", "Avg Lunch")), model.statsDetails.monthAvgFields[2].View()),
		lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%-12s", "Avg OT")), model.statsDetails.monthAvgFields[3].View()),
	)
}

func (model *Model) statsTableView() string {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			model.statsDetails.table.View(),
		  "",
			lipgloss.NewStyle().Bold(true).Underline(true).Render("Selected Month"),
			model.statsMonthAvgView(),
		  "",
		  lipgloss.NewStyle().Bold(true).Underline(true).Render("All Time"),
			model.statsSumView(),
			"",
		)
}


func (model *Model) renderStatisticsView() string {
	var mainDisplay string
	var footer string
	switch model.statsDetails.displayType {
	case graphDisplay:
		footer = "[v] Switch View • [Q] Quit"
		mainDisplay = model.statsGraphView()
	case tableDisplay:
		footer = "[v] Switch View • [Up/Down] Traverse Months • [Q] Quit"
		mainDisplay = model.statsTableView()
	}

	return lipgloss.Place(
		model.width,
		model.height,
		lipgloss.Center,
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Center,
			model.statsHeaderView(),
			"",
			model.statsYearSelectionView(),
			"\n\n",
			mainDisplay,
			"\n",
			lipgloss.NewStyle().Height(2).Render(model.statusLabel),
			lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(footer),
		),
	)
}
