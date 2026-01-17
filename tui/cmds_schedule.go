package tui

import (
	"fmt"
	"fubar/registration"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type offPeriodScheduledMsg struct{}

func (model *Model) scheduleOffPeriod() tea.Cmd {
	return func() tea.Msg {
		var offStart string
		var offEnd string
		offStartInput := model.scheduleInputFields[0].Value()
		if len(offStartInput) == 10 {
			offStart = strings.Replace(offStartInput, "-", "", -1)
		} else if len(offStartInput) == 8 {
			offStart = offStartInput
		} else {
			return errMsg(fmt.Errorf("invalid start date format. must be <YYYY-MM-DD> or <YYYYMMDD>"))
		}
		offEndInput := model.scheduleInputFields[1].Value()
		if len(offEndInput) == 10 {
			offEnd = strings.Replace(offEndInput, "-", "", -1)
		} else if len(offEndInput) == 8 {
			offEnd = offEndInput
		} else {
			return errMsg(fmt.Errorf("invalid end date format. must be <YYYY-MM-DD> or <YYYYMMDD>"))
		}
		offType := model.scheduleInputFields[2].Value()

		err := registration.ScheduleOffPeriod(offStart, offEnd, offType, &model.userConfig)
		if err != nil {
			return errMsg(err)
		}
		return offPeriodScheduledMsg{}
	}
}
