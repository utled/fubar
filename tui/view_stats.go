package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (model *Model) renderStatisticsView() string {
	model.headerFields[0].SetValue("")
	model.headerFields[0].Placeholder = "  [d] Daily"
	model.headerFields[2].SetValue("    Stats")
	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		model.inputStyle.InputField.Render(model.headerFields[0].View()),
		lipgloss.NewStyle().BorderForeground(lipgloss.Color("238")).BorderStyle(lipgloss.NormalBorder()).PaddingLeft(0).Render(model.headerFields[1].View()),
		model.inputStyle.InputField.Render(model.headerFields[2].View()),
		model.inputStyle.InputField.Render(model.headerFields[3].View()),
		model.inputStyle.InputField.Render(model.headerFields[4].View()),
	)

	allSumOne := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%25s", "Worked Days: "), model.statsDetails.allSumFields[idxWorkedDays].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%25s", "Weekdays: "), model.statsDetails.allSumFields[idxWeekdays].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%25s", "Worked Time: "), model.statsDetails.allSumFields[idxWorkedTime].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%25s", "Total OT: "), model.statsDetails.allSumFields[idxTotalOT].View()),
	)
	allSumTwo := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Avg Start: "), model.statsDetails.allSumFields[idxAvgStart].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Avg End: "), model.statsDetails.allSumFields[idxAvgEnd].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Avg Lunch: "), model.statsDetails.allSumFields[idxAvgLunch].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Avg OT: "), model.statsDetails.allSumFields[idxAvgOT].View()),
	)
	allSumThree := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Sic Days: "), model.statsDetails.allSumFields[idxSickDays].View()),
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "Vac Days: "), model.statsDetails.allSumFields[idxVacDays].View()),
		"",
		lipgloss.JoinHorizontal(lipgloss.Left, fmt.Sprintf("%10s", "OT Days: "), model.statsDetails.allSumFields[idxOTDays].View()),
	)
	allSumFields := lipgloss.JoinHorizontal(
		lipgloss.Left,
		allSumOne,
		allSumTwo,
		allSumThree,
		"      ",
	)

	monthSumFields := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%-12s", "Avg Start")), model.statsDetails.monthSumFields[0].View()),
		lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%-12s", "Avg End")), model.statsDetails.monthSumFields[1].View()),
		lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%-12s", "Avg Lunch")), model.statsDetails.monthSumFields[2].View()),
		lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%-12s", "Avg OT")), model.statsDetails.monthSumFields[3].View()),
	)

	yearSelection := lipgloss.JoinHorizontal(
		lipgloss.Center,
		"<-[p] ",
		model.inputStyle.InputField.PaddingRight(1).Render(model.statsDetails.yearSelection.View()),
		" [n]->   ",
	)
	topSpacing := "\n\n"

	var mainDisplay string
	var footer string
	switch model.statsDetails.displayType {
	case graphDisplay:
		footer = "[v] Switch View • [Q] Quit"
		mainDisplay = lipgloss.JoinVertical(
			lipgloss.Center,
			model.statsDetails.graphArea.View(),
			"\n",
			lipgloss.NewStyle().Bold(true).Underline(true).Render("All time"),
			allSumFields,
		)
	case tableDisplay:
		footer = "[v] Switch View • [Up/Down] Traverse Months • [Q] Quit"
		mainDisplay = lipgloss.JoinVertical(
			lipgloss.Center,
			model.statsDetails.table.View(),
			"\n\n\n\n\n\n\n\n\n",
			lipgloss.NewStyle().Bold(true).Underline(true).Render("Selected Month"),
			monthSumFields,
			"\n",
		)
	}

	return lipgloss.Place(
		model.width,
		model.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			header,
			"",
			yearSelection,
			topSpacing,
			mainDisplay,
			"\n",
			model.statusLabel,
			"\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render(footer),
		),
	)
}
