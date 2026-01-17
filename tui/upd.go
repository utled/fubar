package tui

import (
	//"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		model.width = msg.Width
		model.height = msg.Height
		model.statsDetails.graphArea.Height = msg.Height - 30
		return model, nil
	case errMsg:
		model.statusLabel = msg.Error()
		model.dateTable.GotoTop()
		cmds = append(cmds, model.fetchDateState(model.dateTable.SelectedRow()[0]))
		cmds = append(cmds, clearStatusTimer())
		return model, tea.Batch(cmds...)
	case statusMsg:
		cmds = append(cmds, clearStatusTimer())
		return model, tea.Batch(cmds...)
	case clearStatusMsg:
		model.statusLabel = ""
		return model, nil
	case backflushMsg:
		model.statusLabel = "backflush registered successfully"
		cmds = append(cmds, model.fetchDateState(model.dateState.SelectedDate))
		cmds = append(cmds, clearStatusTimer())
		return model, tea.Batch(cmds...)
	case deleteMsg:
		model.confirmationDetails.confirmationMsg = ""
		model.statusLabel = "day deleted successfully"
		cmds = append(cmds, model.fetchDateState(msg.newMaxDate))
		cmds = append(cmds, clearStatusTimer())
		return model, tea.Batch(cmds...)
	case offPeriodScheduledMsg:
		model.statusLabel = "off period scheduled successfully"
		cmds = append(cmds, model.fetchUserConfig())
		cmds = append(cmds, clearStatusTimer())
		return model, tea.Batch(cmds...)
	case offPeriodRemovedMsg:
		model.statusLabel = "off period removed successfully"
		cmds = append(cmds, model.fetchUserConfig())
		cmds = append(cmds, clearStatusTimer())
		return model, tea.Batch(cmds...)
	case configUpdatedMsg:
		model.statusLabel = "user config updated successfully"
		cmds = append(cmds, model.fetchUserConfig())
		cmds = append(cmds, clearStatusTimer())
		return model, tea.Batch(cmds...)
	}

	switch model.state {
	case stateDaily:
		return model.updateDaily(msg)
	case stateBackflush:
		return model.updateBackflush(msg)
	case stateStatistics:
		return model.updateStats(msg)
	case stateSchedule:
		return model.updateSchedule(msg)
	case stateConfig:
		return model.updateConfig(msg)
	case stateConfirm:
		return model.updateConfirmation(msg)
	}
	return model, nil
}
