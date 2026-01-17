package tui

import (
	"fubar/registration"

	tea "github.com/charmbracelet/bubbletea"
)

type backflushMsg struct {
}

func (model *Model) registerBackflush() tea.Cmd {
	return func() tea.Msg {
		err := registration.RegisterBackflush(model.backflushInputFields.Value(), &model.dateState, &model.userConfig)
		if err != nil {
			return errMsg(err)
		}
		return backflushMsg{}
	}
}
