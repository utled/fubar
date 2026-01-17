package tui

import tea "github.com/charmbracelet/bubbletea"

func (model *Model) updateConfirmation(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		}
		switch msg.Type {
		case tea.KeyEnter:
			switch model.confirmationDetails.confirmationType {
			case deleteDate:
				model.state = stateDaily
				return model, model.deleteDate()
			case deleteSchedule:
				model.state = model.prevState
				return model, model.removeOffPeriod()
			}
		case tea.KeyEsc:
			switch model.confirmationDetails.confirmationType {
			case deleteDate:
				model.state = stateDaily
			case deleteSchedule:
				model.state = stateSchedule
				return model, nil
			}
		}
	}
	return model, nil
}
