package tui

func (model *Model) View() string {
	dailyView := model.renderDailyView()
	if model.state == stateBackflush {
		backflushView := model.renderBackflushView()
		return backflushView
	}
	if model.state == stateStatistics {
		statisticsView := model.renderStatisticsView()
		return statisticsView
	}
	if model.state == stateSchedule {
		scheduleView := model.renderScheduleView()
		return scheduleView
	}
	if model.state == stateConfig {
		configView := model.renderConfigView()
		return configView
	}
	if model.state == stateConfirm {
		confirmationView := model.renderConfirmationView()
		return confirmationView
	}
	return dailyView
}
