package tui

import (
	"fmt"
	"fubar/helpers"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

func (model *Model) renderDailyView() string {
	model.headerFields[0].SetValue("    Daily")
	model.headerFields[2].Placeholder = "  [x] Stats"
	model.headerFields[2].SetValue("")
	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		model.inputStyle.InputField.Render(model.headerFields[0].View()),
		lipgloss.NewStyle().BorderForeground(lipgloss.Color("238")).BorderStyle(lipgloss.NormalBorder()).PaddingLeft(0).Render(model.headerFields[1].View()),
		model.inputStyle.InputField.Render(model.headerFields[2].View()),
		model.inputStyle.InputField.Render(model.headerFields[3].View()),
		model.inputStyle.InputField.Render(model.headerFields[4].View()),
	)

	subHeaderOne := lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Registration up to date: ",
		strconv.FormatBool(model.dateState.ReportUpToDate),
	)

	var subHeaderTwo string
	if !model.dateState.ReportUpToDate {
		subHeaderTwo = lipgloss.JoinHorizontal(
			lipgloss.Center,
			"Last registered date: ",
			model.dateState.MaxCompletedDate,
		)
	} else {
		subHeaderTwo = ""
	}

	subHeaderThree := lipgloss.JoinHorizontal(
		lipgloss.Center,
		model.dateState.SelectedDate,
	)

	subHeaderFour := lipgloss.JoinHorizontal(
		lipgloss.Center,
		"Balance: ",
		helpers.DecimalToTime(model.dateState.TotalBalance),
	)

	rowOne := lipgloss.JoinHorizontal(
		lipgloss.Center,
		fmt.Sprintf("%16s", "[s] Start: "),
		model.inputStyle.InputField.Render(model.dailyInputFields[idxStart].View()),
		fmt.Sprintf("%16s", "[t] Type: "),
		model.inputStyle.InputField.Render(model.dailyInputFields[idxDayType].View()),
		fmt.Sprintf("%7s", ""),
	)
	rowTwo := lipgloss.JoinHorizontal(
		lipgloss.Center,
		fmt.Sprintf("%16s", "[l] Lunch: "),
		model.inputStyle.InputField.Render(model.dailyInputFields[idxLunch].View()),
		fmt.Sprintf("%16s", "[o] Overtime: "),
		model.inputStyle.InputField.Render(model.dailyInputFields[idxOvertime].View()),
		fmt.Sprintf("%7s", ""),
	)
	rowThree := lipgloss.JoinHorizontal(
		lipgloss.Center,
		fmt.Sprintf("%16s", "[e] End: "),
		model.inputStyle.InputField.Render(model.dailyInputFields[idxEnd].View()),
		fmt.Sprintf("%16s", "Day Total: "),
		model.inputStyle.InputField.Render(model.dailyInputFields[idxDayTotal].View()),
		fmt.Sprintf("%7s", ""),
	)
	rowFour := lipgloss.JoinHorizontal(
		lipgloss.Center,
		fmt.Sprintf("%16s", "[a] Additional: "),
		model.inputStyle.InputField.Render(model.dailyInputFields[idxAdditional].View()),
		fmt.Sprintf("%16s", "Day Balance: "),
		model.inputStyle.InputField.Render(model.dailyInputFields[idxDayBalance].View()),
		fmt.Sprintf("%7s", ""),
	)

	footerOne := "[Up/Down] Traverse Days • [Left/Right] Switch Input"
	footerTwo := "[Enter] Save • [Esc] Revert Changes • [Q] Quit"

	return lipgloss.Place(
		model.width,
		model.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			subHeaderOne,
			subHeaderTwo,
			"\n\n",
			lipgloss.NewStyle().Bold(true).Underline(true).Render(subHeaderThree),
			subHeaderFour,
			"\n\n",
			rowOne,
			rowTwo,
			rowThree,
			rowFour,
			"\n",
			model.dateTable.View(),
			"\n ",
			model.statusLabel,
			"\n\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(footerOne),
			lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(footerTwo),
		),
	)
}
