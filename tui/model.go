package tui

import (
	"fubar/data"
	"fubar/utils"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	state                sessionState
	prevState            sessionState
	width                int
	height               int
	inputStyle           styles
	userConfig           data.UserConfig
	statusLabel          string
	dateState            data.ReportState
	headerFields         []textinput.Model
	dailyInputFields     []textinput.Model
	dailyInputFocus      int
	timesheet            []*data.WorkDateRecord
	dateTable            table.Model
	backflushInputFields textinput.Model
	statsDetails         statsDetails
	scheduleInputFields  []textinput.Model
	configInputFields    []textinput.Model
	confirmationDetails  confirmationDetails
	otherInputFocus      int
}

type styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func inputStyles() styles {
	style := new(styles)
	style.BorderColor = lipgloss.Color("238")
	style.InputField = lipgloss.NewStyle().
		BorderForeground(style.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		PaddingLeft(2)
	return *style
}

func tableStyles() table.Styles {
	tableStyle := table.DefaultStyles()

	tableStyle.Header = tableStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("238")).
		BorderBottom(true).
		Bold(true)

	tableStyle.Selected = tableStyle.Header.
		BorderForeground(lipgloss.Color("238")).
		BorderBottom(true).
		Bold(true)

	tableStyle.Selected = tableStyle.Selected.
		Background(lipgloss.Color("234")).
		Bold(true).
		PaddingLeft(0)

	return tableStyle
}

func NewModel() Model {
	//                 //
	//                 //
	// Daily Model     //
	//                 //
	//                 //
	inputStyle := inputStyles()
	tableStyle := tableStyles()

	headerFields := make([]textinput.Model, 5)
	for i := range headerFields {
		inputField := textinput.New()
		inputField.Width = 15
		inputField.Prompt = ""
		switch i {
		case 0:
			inputField.SetValue("    Daily")
		case 1:
			inputField.Placeholder = "  [b] Backflush"
		case 2:
			inputField.Placeholder = "  [x] Stats"
		case 3:
			inputField.Placeholder = "[y] Schedule"
		case 4:
			inputField.Placeholder = " [z] Config"
		}
		headerFields[i] = inputField
	}

	dailyInputFields := make([]textinput.Model, 8)
	for i := range dailyInputFields {
		inputField := textinput.New()
		inputField.Width = 7
		inputField.Prompt = ""
		switch i {
		case idxStart:
			inputField.CharLimit = 5
		case idxLunch:
			inputField.CharLimit = 3
		case idxEnd:
			inputField.CharLimit = 5
		case idxAdditional:
			inputField.CharLimit = 3
		}
		dailyInputFields[i] = inputField
	}
	dailyInputFields[idxStart].Focus()

	columns := []table.Column{
		{Title: "Date", Width: 12},
		{Title: "Type", Width: 7},
		{Title: "Start", Width: 7},
		{Title: "Lunch", Width: 7},
		{Title: "End", Width: 7},
		{Title: "+Time", Width: 7},
		{Title: "Day Tot", Width: 7},
		{Title: "OT", Width: 7},
		{Title: "Day +/-", Width: 7},
		{Title: "Total +/-", Width: 9},
	}

	dateTable := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	dateTable.SetStyles(tableStyle)

	//                 //
	//                 //
	// Backflush model
	//                 //
	//                 //
	backflushInputField := textinput.New()
	backflushInputField.Width = 20
	backflushInputField.CharLimit = 4

	//                 //
	//                 //
	// Stats model
	//                 //
	//                 //
	yearSelection := textinput.New()
	yearSelection.Prompt = ""
	yearSelection.Width = 4
	yearSelection.CharLimit = 4

	allSumFields := make([]textinput.Model, 11)
	for i := range allSumFields {
		inputField := textinput.New()
		inputField.Prompt = ""
		inputField.Width = 10
		inputField.CharLimit = 10
		allSumFields[i] = inputField
	}

	monthSumFields := make([]textinput.Model, 4)
	for i := range monthSumFields {
		inputField := textinput.New()
		inputField.Prompt = ""
		inputField.Width = 12
		inputField.CharLimit = 12
		monthSumFields[i] = inputField
	}

	graphArea := viewport.New(55, 20)
	graphArea.SetContent("Loading graph...")

	statsColumns := []table.Column{
		{Title: "Month", Width: 9},
		{Title: "Week", Width: 5},
		{Title: "Wrkd", Width: 5},
		{Title: "WrkdT", Width: 8},
		{Title: "Vac", Width: 4},
		{Title: "Sic", Width: 4},
		{Title: "Wknd", Width: 4},
		{Title: "Off", Width: 4},
		{Title: "OT", Width: 4},
		{Title: "OTT", Width: 5},
	}

	statsTable := table.New(
		table.WithColumns(statsColumns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	statsTable.SetStyles(tableStyle)

	statsDetails := statsDetails{
		yearSelection:  yearSelection,
		table:          statsTable,
		graphArea:      graphArea,
		allSumFields:   allSumFields,
		monthSumFields: monthSumFields,
	}

	//                 //
	//                 //
	// Schedule model  //
	//                 //
	//                 //
	scheduleInputFields := make([]textinput.Model, 3)
	for i := range scheduleInputFields {
		inputField := textinput.New()
		inputField.Width = 20
		switch i {
		case 0:
			inputField.Placeholder = "<YYYY-MM-DD/YYYYMMDD>"
		case 1:
			inputField.Placeholder = "<YYYY-MM-DD/YYYYMMDD>"
		}

		inputField.CharLimit = 10
		scheduleInputFields[i] = inputField
	}

	//                 //
	//                 //
	// Config model    //
	//                 //
	//                 //
	configInputFields := make([]textinput.Model, 2)
	for i := range configInputFields {
		inputField := textinput.New()
		inputField.Width = 12
		switch i {
		case 0:
			inputField.CharLimit = 2
			inputField.Placeholder = "<MM>"
		case 1:
			inputField.CharLimit = 5
			inputField.Placeholder = "<HH:MM/HHMM>"
		}

		configInputFields[i] = inputField
	}

	return Model{
		state:                stateDaily,
		headerFields:         headerFields,
		dailyInputFields:     dailyInputFields,
		backflushInputFields: backflushInputField,
		statsDetails:         statsDetails,
		scheduleInputFields:  scheduleInputFields,
		configInputFields:    configInputFields,
		dateTable:            dateTable,
		inputStyle:           inputStyle,
	}
}

func (model *Model) Init() tea.Cmd {
	model.dateState.SelectedDate = time.Now().Format(utils.DateLayout)
	return model.fetchUserConfig()
}
