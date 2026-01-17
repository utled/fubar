package tui

import "github.com/charmbracelet/lipgloss"

func (model *Model) renderConfirmationView() string {
	modalStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238")).
		Padding(0).
		Width(40)

	var confirmationHeader string
	switch model.confirmationDetails.confirmationType {
	case deleteDate:
		confirmationHeader = "Delete Date"
	case deleteSchedule:
		confirmationHeader = "Delete Scheduled off period"
	}

	form := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Bold(true).Render(confirmationHeader),
		"\n",
		lipgloss.NewStyle().Align(lipgloss.Center).Bold(false).Italic(true).Render(model.confirmationDetails.confirmationMsg),
		"\n",
		lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render("Enter: Confirm • Esc: Cancel"),
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
