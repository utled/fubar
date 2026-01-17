package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func Launch() {
	model := NewModel()
	program := tea.NewProgram(&model, tea.WithAltScreen())
	_, err := program.Run()
	if err != nil {
		log.Fatalf("failed to run bubble program: %v", err)
	}
}
