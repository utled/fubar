package tui

import tea "github.com/charmbracelet/bubbletea"

func (model *Model) updateBackflush(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "t":
			nextDayType := ParseDayType(model.backflushInputFields.Value()).Next().String()
			model.backflushInputFields.SetValue(nextDayType)
			return model, nil
		case "ctrl+c", "q":
			return model, tea.Quit
		}
		switch msg.Type {
		case tea.KeyEnter:
			model.state = model.prevState
			return model, model.registerBackflush()
		case tea.KeyEsc:
			model.state = model.prevState
			return model, nil
		}
	}
	var cmd tea.Cmd
	model.backflushInputFields, cmd = model.backflushInputFields.Update(msg)
	return model, cmd
}
