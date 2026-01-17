package tui

import (
	"fubar/helpers"
	"fubar/registration"

	tea "github.com/charmbracelet/bubbletea"
)

type deleteMsg struct {
	newMaxDate string
}

func (model *Model) deleteDate() tea.Cmd {
	return func() tea.Msg {
		err := registration.DeleteDate(&model.dateState)
		if err != nil {
			return errMsg(err)
		}
		newMaxDate, err := helpers.PreviousDateFromString(model.dateState.SelectedDate)
		if err != nil {
			return errMsg(err)
		}
		return deleteMsg{newMaxDate: newMaxDate}
	}
}

type offPeriodRemovedMsg struct{}

func (model *Model) removeOffPeriod() tea.Cmd {
	return func() tea.Msg {
		err := registration.RemoveScheduledOffPeriod()
		if err != nil {
			return errMsg(err)
		}
		return offPeriodRemovedMsg{}
	}
}
