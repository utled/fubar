package tui

import tea "github.com/charmbracelet/bubbletea"

func (model *Model) updateConfig(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		}
		switch msg.Type {
		case tea.KeyLeft, tea.KeyUp, tea.KeyRight, tea.KeyDown, tea.KeyTab:
			model.configInputFields[model.otherInputFocus].Blur()
			if msg.Type == tea.KeyLeft || msg.Type == tea.KeyUp {
				model.otherInputFocus = (model.otherInputFocus - 1 + len(model.configInputFields)) % len(model.configInputFields)
			} else {
				model.otherInputFocus = (model.otherInputFocus + 1) % len(model.configInputFields)
			}
			return model, model.configInputFields[model.otherInputFocus].Focus()
		case tea.KeyEnter:
			model.state = model.prevState
			return model, model.updateUserConfig()
		case tea.KeyEsc:
			var cmd tea.Cmd
			model.syncInputs()
			model.state = model.prevState
			model.configInputFields[model.otherInputFocus], cmd = model.configInputFields[model.otherInputFocus].Update(msg)
			return model, cmd
		}
	}
	var cmd tea.Cmd
	model.configInputFields[model.otherInputFocus], cmd = model.configInputFields[model.otherInputFocus].Update(msg)
	return model, cmd
}
