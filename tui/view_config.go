package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (model *Model) renderConfigView() string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(0).
		Width(40)

	lunchRow := lipgloss.JoinHorizontal(
		lipgloss.Left,
		fmt.Sprintf("%15s", "Default Lunch: "),
		model.configInputFields[0].View(),
	)
	dayLengthRow := lipgloss.JoinHorizontal(
		lipgloss.Left,
		fmt.Sprintf("%15s", "Day Length: "),
		model.configInputFields[1].View(),
	)

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("User Configuration"),
		"",
		lunchRow,
		dayLengthRow,
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render("Enter: Save • Esc: Cancel"),
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
