package tui

import "github.com/charmbracelet/lipgloss"

func (model *Model) renderBackflushView() string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(0).
		Width(40)

	inputRow := lipgloss.JoinHorizontal(
		lipgloss.Left,
		"[t] Day Type: ",
		model.backflushInputFields.View(),
	)

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("Register backflush"),
		lipgloss.NewStyle().Bold(false).Italic(true).Render("Autofills all dates back to last registered date with the selected day type"),
		"",
		inputRow,
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
