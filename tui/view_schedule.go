package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (model *Model) renderScheduleView() string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(0).
		Width(40)

	startRow := lipgloss.JoinHorizontal(
		lipgloss.Left,
		fmt.Sprintf("%15s", "Start Date: "),
		model.scheduleInputFields[0].View(),
	)
	endRow := lipgloss.JoinHorizontal(
		lipgloss.Left,
		fmt.Sprintf("%15s", "End Date: "),
		model.scheduleInputFields[1].View(),
	)
	typeRow := lipgloss.JoinHorizontal(
		lipgloss.Left,
		fmt.Sprintf("%15s", "[t] Day Type: "),
		model.scheduleInputFields[2].View(),
	)

	form := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Bold(true).Render("Schedule off period"),
		"",
		startRow,
		endRow,
		typeRow,
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render("Enter: Save • Del: Delete • Esc: Cancel"),
	)

	return lipgloss.Place(
		model.width,
		model.height,
		lipgloss.Center,
		lipgloss.Center,
		modalStyle.Render(form),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("238")),
	)
}
