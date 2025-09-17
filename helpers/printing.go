package helpers

import (
	"fTime/data"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func PrintHeader(withSupportText bool, state *data.ReportState) {
	fmt.Println()
	fmt.Printf("%72s", "┏━━━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━━┓\n")
	fmt.Printf("%72s", "┃┏┓┏┓┃┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┓┏┓┃\n")
	fmt.Printf("%72s", "┗┛┃┃┗┛┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃┃┃\n")
	fmt.Printf("%72s", "━━┃┃━━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃┃┃\n")
	fmt.Printf("%72s", "━┏┛┗┓━┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┛┗┛┃\n")
	fmt.Printf("%72s", "━┗━━┛━━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━━┛\n")
	fmt.Printf("%72s", "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("%72s", "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	if withSupportText {
		fmt.Printf("%72s", "cmd -> Display available commands           Ctrl+C -> Close program\n")
		fmt.Println()
	}
	if !state.ReportUpToDate {
		fmt.Printf("%4s", "")
		fmt.Printf("%-67s", "There are missing regitrations")
		fmt.Println()
		fmt.Printf("%4s", "")
		fmt.Printf("%-15s%-10s", "Last completed date:", state.MaxCompletedDate)
		fmt.Println()
		fmt.Println()
	} else {
		fmt.Printf("%4s", "")
		fmt.Printf("%-67s", "Registrations are up to date")
		fmt.Println()
		fmt.Println()
	}
	fmt.Printf("%4s", "")
	fmt.Printf("%-15s%-10s", "Selected date:", state.SelectedDate)
	fmt.Println()
	fmt.Println()

}

func PrintSelectedDate(state *data.ReportState) {
	err := ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}
	PrintHeader(true, state)
	fmt.Printf("%-12s", "Start: ")
	fmt.Printf("%-20s", state.SelectedRecord.StartTime.String)
	fmt.Printf("%-15s", "Day Total: ")
	fmt.Printf("%-20s", state.SelectedRecord.DayTotal.String)
	fmt.Printf("%-15s", "Projected End: ")
	fmt.Printf("%-20s\n", state.ProjectedEnd)

	fmt.Printf("%-12s", "Lunch: ")
	fmt.Printf("%-20d", state.SelectedRecord.LunchDuration.Int16)
	fmt.Printf("%-15s", "Day Balance: ")
	fmt.Printf("%-20.2f\n", state.SelectedRecord.DayBalance.Float64)

	fmt.Printf("%-12s", "End: ")
	fmt.Printf("%-20s", state.SelectedRecord.EndTime.String)
	fmt.Printf("%-15s", "Total Balance: ")
	fmt.Printf("%-20.2f\n", state.SelectedRecord.MovingBalance.Float64)

	fmt.Printf("%-12s", "Additional: ")
	fmt.Printf("%-20d\n", state.SelectedRecord.AdditionalTime.Int16)
	fmt.Printf("%-12s", "Overtime: ")
	fmt.Printf("%-20s\n", fmt.Sprintf("%t", state.SelectedRecord.Overtime.Bool))
	fmt.Printf("%-12s", "Type: ")
	fmt.Printf("%-20s\n", state.SelectedRecord.DayType.String)

}

func PrintCommands(state *data.ReportState) {
	err := ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}
	PrintHeader(false, state)
	availableCommands := []string{
		"\n_____DISPLAY_____________________________________________________________________________________________",
		"today[t]                                       -> Display current date",
		"switch[sw] <YYYYMMDD>                          -> Switch display date",
		"last[la] <INT(days)>                           -> Display last X days as table",
		"range[ra] <YYYYMMDD> <YYYYMMD>                 -> Display date range as table",
		"clear[c]                                       -> Clear terminal",
		"\n_____DAILY ACTIONS_______________________________________________________________________________________",
		"start[s] <MMSS>                                -> Set/modify start time for displayed date",
		"end[e] <MMSS> [optional]<ot/-ot>>              -> Set/modify end time for displayed date",
		"ot/-ot                                         -> Set/remove excess time as overtime for displayed date",
		"lunch[l] <INT(minutes)>                        -> Set/modify lunch time for displayed date",
		"addit[ad] <INT(minutes)>                       -> Set/modify additional time for displayed date>",
		"off/vac/sic                                    -> Flag displayed date as non working day/partial day",
		"norm                                           -> Revert an off day back to normal",
		"\n_____SCHEDULING__________________________________________________________________________________________",
		"sched[sc] <YYYYMMDD> <YYYYMMDD> <off/vac/sic>  -> Schedule date period for coming off time",
		"sched[sc] remove                               -> Remove scheduled period",
		"back <norm/off/vac/sic>                        -> Backfill all non registered days back to last completed date",
		"\n_____DEFAULT CONFIGURATIONS______________________________________________________________________________",
		"conflunch <INT(minutes)>                       -> Update default lunch duration",
		"conflength <MMSS>                              -> Update default length of day",
		"\n_____GENERAL_____________________________________________________________________________________________",
		"cmd                                            -> Display available commands",
		"Ctrl+C                                         -> Close program",
		"\n_____OTHER_______________________________________________________________________________________________",
		"delete[dl]                                     -> Delete selected date (can only delete last registered date)",
	}

	for _, command := range availableCommands {
		fmt.Println(command)
	}
	fmt.Println()
}

func PrintDateRange(dateRange []*data.WorkDateRecord, reversed bool, state *data.ReportState) {
	err := ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}
	PrintHeader(false, state)

	fmt.Printf("%-15s", "Date")
	fmt.Printf("%-15s", "Type")
	fmt.Printf("%-15s", "Start")
	fmt.Printf("%-15s", "Lunch")
	fmt.Printf("%-15s", "End")
	fmt.Printf("%-15s", "Additional")
	fmt.Printf("%-15s", "Day Total")
	fmt.Printf("%-15s", "Overtime")
	fmt.Printf("%-15s", "Day Balance")
	fmt.Printf("%-15s\n", "Total Balance")
	for i := 1; i < 149; i++ {
		fmt.Print("_")
	}
	fmt.Println()

	if !reversed {
		for index := len(dateRange) - 1; index >= 0; index-- {
			fmt.Printf("%-15s", dateRange[index].WorkDate)
			fmt.Printf("%-15s", dateRange[index].DayType.String)
			fmt.Printf("%-15s", dateRange[index].StartTime.String)
			fmt.Printf("%-15d", dateRange[index].LunchDuration.Int16)
			fmt.Printf("%-15s", dateRange[index].EndTime.String)
			fmt.Printf("%-15d", dateRange[index].AdditionalTime.Int16)
			fmt.Printf("%-15s", dateRange[index].DayTotal.String)
			fmt.Printf("%-15s", fmt.Sprintf("%t", dateRange[index].Overtime.Bool))
			fmt.Printf("%-15.2f", dateRange[index].DayBalance.Float64)
			fmt.Printf("%-15.2f\n", dateRange[index].MovingBalance.Float64)
		}
		fmt.Println()
	} else {
		for _, date := range dateRange {
			fmt.Printf("%-15s", date.WorkDate)
			fmt.Printf("%-15s", date.DayType.String)
			fmt.Printf("%-15s", date.StartTime.String)
			fmt.Printf("%-15d", date.LunchDuration.Int16)
			fmt.Printf("%-15s", date.EndTime.String)
			fmt.Printf("%-15d", date.AdditionalTime.Int16)
			fmt.Printf("%-15s", date.DayTotal.String)
			fmt.Printf("%-15s", fmt.Sprintf("%t", date.Overtime.Bool))
			fmt.Printf("%-15.2f", date.DayBalance.Float64)
			fmt.Printf("%-15.2f\n", date.MovingBalance.Float64)
		}
	}

}

var clearFunctions map[string]func()

func InitClearFunctions() {
	clearFunctions = make(map[string]func())
	clearFunctions["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearFunctions["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ClearTerminal() error {
	clearFunction, ok := clearFunctions[runtime.GOOS]
	if ok {
		clearFunction()
	} else {
		return fmt.Errorf("unable to clear the terminal")
	}
	return nil
}
