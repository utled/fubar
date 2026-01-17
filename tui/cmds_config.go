package tui

import (
	"fmt"
	"fubar/data"
	"fubar/helpers"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

type configUpdatedMsg struct{}

func (model *Model) updateUserConfig() tea.Cmd {
	return func() tea.Msg {
		defaultLunch, err := strconv.Atoi(model.configInputFields[0].Value())
		if err != nil {
			return errMsg(err)
		}
		if defaultLunch < 0 {
			return errMsg(fmt.Errorf("lunch Duration can't be a negative value"))
		}
		err = data.UpdateDefaultLunch(defaultLunch)
		if err != nil {
			return errMsg(err)
		}

		defaultDayLength, err := helpers.FormatValidTimeString(model.configInputFields[1].Value())
		err = data.UpdateDefaultLength(defaultDayLength)
		if err != nil {
			return errMsg(err)
		}
		return configUpdatedMsg{}
	}
}
