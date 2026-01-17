package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error
type statusMsg struct{}

func (model *Model) setStatusMsg(msg string) tea.Cmd {
	return func() tea.Msg {
		model.statusLabel = msg
		return statusMsg{}
	}
}

type clearStatusMsg struct{}

func clearStatusTimer() tea.Cmd {
	return tea.Tick(time.Second*4, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}
