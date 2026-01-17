package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) updateSchedule(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "t":
			nextDayType := ParseDayType(model.scheduleInputFields[2].Value()).Next().String()
			model.scheduleInputFields[2].SetValue(nextDayType)
			return model, nil
		case "ctrl+c", "q":
			return model, tea.Quit
		}
		switch msg.Type {
		case tea.KeyLeft, tea.KeyUp, tea.KeyRight, tea.KeyDown, tea.KeyTab:
			focusCycleLimit := 2
			model.scheduleInputFields[model.otherInputFocus].Blur()
			if msg.Type == tea.KeyLeft || msg.Type == tea.KeyUp {
				model.otherInputFocus = (model.otherInputFocus - 1 + focusCycleLimit) % focusCycleLimit
			} else {
				model.otherInputFocus = (model.otherInputFocus + 1) % focusCycleLimit
			}
			return model, model.scheduleInputFields[model.otherInputFocus].Focus()
		case tea.KeyEnter:
			model.state = model.prevState
			return model, model.scheduleOffPeriod()
		case tea.KeyDelete:
			if model.userConfig.OffStart.String != "" {
				model.confirmationDetails.confirmationType = deleteSchedule
				model.confirmationDetails.confirmationMsg = fmt.Sprint(
					"Please confirm deletion of period:\n",
					model.userConfig.OffStart.String, " to ", model.userConfig.OffEnd.String,
				)
				model.state = stateConfirm
				return model, nil
			}
		case tea.KeyEsc:
			var cmd tea.Cmd
			model.syncInputs()
			model.state = model.prevState
			model.scheduleInputFields[model.otherInputFocus], cmd = model.scheduleInputFields[model.otherInputFocus].Update(msg)
			return model, cmd
		}
	}
	var cmd tea.Cmd
	model.scheduleInputFields[model.otherInputFocus], cmd = model.scheduleInputFields[model.otherInputFocus].Update(msg)
	return model, cmd
}
