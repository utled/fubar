package helpers

import (
	"fTime/data"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// PrintHeader clears the terminal and prints the header and subheader.
// To be used before any other printing functions.
func PrintHeader(state *data.ReportState) {
	err := ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()

	fmt.Printf("%65s", "    ██████  ███████████  ███                          \n")
	fmt.Printf("%65s", "   ███░░███░█░░░███░░░█ ░░░                           \n")
	fmt.Printf("%65s", "  ░███ ░░░ ░   ░███  ░  ████  █████████████    ██████ \n")
	fmt.Printf("%65s", " ███████       ░███    ░░███ ░░███░░███░░███  ███░░███\n")
	fmt.Printf("%65s", "░░░███░        ░███     ░███  ░███ ░███ ░███ ░███████ \n")
	fmt.Printf("%65s", "  ░███         ░███     ░███  ░███ ░███ ░███ ░███░░░  \n")
	fmt.Printf("%65s", "  █████        █████    █████ █████░███ █████░░██████ \n")
	fmt.Printf("%65s", " ░░░░░        ░░░░░    ░░░░░ ░░░░░ ░░░ ░░░░░  ░░░░░░  \n")
	fmt.Printf("%65s", "                                                      \n")
	fmt.Printf("%72s", "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	fmt.Printf("%72s", "cmd -> Display available commands           Ctrl+C -> Close program\n")
	fmt.Println()

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
	fmt.Printf("%72s", "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Println()
	fmt.Println()

}

func PrintSelectedDate(state *data.ReportState) {
	PrintHeader(state)

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
	fmt.Printf("%-20.2f\n", state.SelectedRecord.TotalBalance.Float64)

	fmt.Printf("%-12s", "Additional: ")
	fmt.Printf("%-20d\n", state.SelectedRecord.AdditionalTime.Int16)
	fmt.Printf("%-12s", "Overtime: ")
	fmt.Printf("%-20s\n", fmt.Sprintf("%t", state.SelectedRecord.Overtime.Bool))
	fmt.Printf("%-12s", "Type: ")
	fmt.Printf("%-20s\n", state.SelectedRecord.DayType.String)

	fmt.Println()

}

func PrintCommands(state *data.ReportState) {
	PrintHeader(state)
	availableCommands := []string{
		"\n_____DISPLAY_____________________________________________________________________________________________",
		"today[t]                                       -> Display current date",
		"next[n]                                        -> Display next date",
		"prev[p]                                        -> Display previous date",
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
		"sched[sc] show                                 -> Display scheduled period",
		"back <norm/off/vac/sic>                        -> Backfill all non registered days back to last completed date",
		"\n_____DEFAULT CONFIGURATIONS______________________________________________________________________________",
		"conf lunch[l] <INT(minutes)>                   -> Update default lunch duration",
		"conf length[le] <MMSS>                         -> Update default length of day",
		"conf show[s]                                   -> Display current user defaults",
		"\n_____GENERAL_____________________________________________________________________________________________",
		"cmd                                            -> Display available commands",
		"Ctrl+C                                         -> Close program",
		"\n_____OTHER_______________________________________________________________________________________________",
		"delete[dl]                                     -> Delete selected date (can only delete last registered date)",
		"\n_____STATISTICS__________________________________________________________________________________________",
		"stats[st] all[a]                               -> Display all time statistics",
		"stats[st] sum[s] <YYYY>                        -> Display yearly summary grouped by month",
		"stats[st] year[y] <YYYY>                       -> Display statistics for specified year",
		"stats[st] month[m] <INT(monthnum)> <YYYY>      -> Display statistics for specified month",
		"stats[st] days[d] <INT(days)>                  -> Display statistics for specified num of days back",
		"stats[st] range[r] <YYYYMMDD> <YYYYMMDD>       -> Display statistics for specified range of dates",
	}

	for _, command := range availableCommands {
		fmt.Println(command)
	}
	fmt.Println()
}

func PrintScheduledOffPeriod(userConfig *data.UserConfig, state *data.ReportState) {
	PrintHeader(state)
	if userConfig.OffStart.String == "" {
		fmt.Println("No off period is scheduled")
	} else {
		fmt.Printf("%-30s", "Scheduled off period, Start: ")
		fmt.Print(userConfig.OffStart.String, "\n")
		fmt.Printf("%-30s", "Scheduled off period, End: ")
		fmt.Print(userConfig.OffEnd.String, "\n")
		fmt.Printf("%-30s", "Scheduled off period, Type: ")
		fmt.Print(userConfig.OffType.String, "\n")
	}
	fmt.Println()
}

func PrintUserConfig(userConfig *data.UserConfig, state *data.ReportState) {
	PrintHeader(state)
	fmt.Printf("%-20s", "Default Lunch: ")
	fmt.Print(userConfig.DefaultLunch.Int16, "\n")
	fmt.Printf("%-20s", "Default Day Length: ")
	fmt.Print(userConfig.DefaultDayLength.String, "\n")
	fmt.Println()
}

func PrintDateRange(dateRange []*data.WorkDateRecord, ascending bool, state *data.ReportState) {
	PrintHeader(state)

	fmt.Printf("%-12s", "Date")
	fmt.Printf("%-6s", "Type")
	fmt.Printf("%-10s", "Start")
	fmt.Printf("%-7s", "Lunch")
	fmt.Printf("%-10s", "End")
	fmt.Printf("%-7s", "+Time")
	fmt.Printf("%-11s", "Day Total")
	fmt.Printf("%-10s", "Overtime")
	fmt.Printf("%-9s", "Day +/-")
	fmt.Printf("%-10s\n", "Total +/-")
	for i := 1; i < 92; i++ {
		fmt.Print("_")
	}
	fmt.Println()

	if ascending {
		for _, date := range dateRange {
			fmt.Printf("%-12s", date.WorkDate)
			fmt.Printf("%-6s", date.DayType.String)
			fmt.Printf("%-10s", date.StartTime.String)
			fmt.Printf("%-7d", date.LunchDuration.Int16)
			fmt.Printf("%-10s", date.EndTime.String)
			fmt.Printf("%-7d", date.AdditionalTime.Int16)
			fmt.Printf("%-11s", date.DayTotal.String)
			fmt.Printf("%-10s", fmt.Sprintf("%t", date.Overtime.Bool))
			fmt.Printf("%-9.2f", date.DayBalance.Float64)
			fmt.Printf("%-10.2f\n", date.TotalBalance.Float64)
		}
	} else {
		for index := len(dateRange) - 1; index >= 0; index-- {
			fmt.Printf("%-12s", dateRange[index].WorkDate)
			fmt.Printf("%-6s", dateRange[index].DayType.String)
			fmt.Printf("%-10s", dateRange[index].StartTime.String)
			fmt.Printf("%-7d", dateRange[index].LunchDuration.Int16)
			fmt.Printf("%-10s", dateRange[index].EndTime.String)
			fmt.Printf("%-7d", dateRange[index].AdditionalTime.Int16)
			fmt.Printf("%-11s", dateRange[index].DayTotal.String)
			fmt.Printf("%-10s", fmt.Sprintf("%t", dateRange[index].Overtime.Bool))
			fmt.Printf("%-9.2f", dateRange[index].DayBalance.Float64)
			fmt.Printf("%-10.2f\n", dateRange[index].TotalBalance.Float64)
		}
	}

	fmt.Println()

}

func PrintMonthlySummary(monthlySummary []*data.MonthStats, title string, state *data.ReportState) {
	PrintHeader(state)

	fmt.Print("\n", title, "\n\n")
	PrintWorkedDaysDiagram(monthlySummary)
	fmt.Println()
	for i := 1; i < 112; i++ {
		fmt.Print("_")
	}
	fmt.Print("\n\n")

	fmt.Printf("%-11s", "Month")
	fmt.Printf("%-10s", "Weekdays")
	fmt.Printf("%-13s", "Worked Days")
	fmt.Printf("%-13s", "Worked Time")
	fmt.Printf("%-15s", "Vacation Days")
	fmt.Printf("%-11s", "Sick Days")
	fmt.Printf("%-11s", "Wknd Days")
	fmt.Printf("%-10s", "Off Days")
	fmt.Printf("%-9s", "OT Days")
	fmt.Printf("%-10s\n", "Total OT")
	for i := 1; i < 112; i++ {
		fmt.Print("_")
	}
	fmt.Println()

	for _, month := range monthlySummary {
		fmt.Printf("%-11s", month.Month)
		fmt.Printf("%-10d", month.TotalWeekDays)
		fmt.Printf("%-13d", month.WorkedDays)
		fmt.Printf("%-13s", month.WorkedTime)
		fmt.Printf("%-15d", month.VacationDays)
		fmt.Printf("%-11d", month.SickDays)
		fmt.Printf("%-11d", month.WeekendDays)
		fmt.Printf("%-10d", month.OffDays)
		fmt.Printf("%-9d", month.OverTimeDays)
		fmt.Printf("%-10.2f\n", month.TotalOvertime.Float64)
	}

	fmt.Println()
}

func PrintFullStatistics(fullStatistics *data.FullStats, title string, state *data.ReportState) {
	PrintHeader(state)

	fmt.Print("\n", title, "\n")
	for i := 1; i < 93; i++ {
		fmt.Print("_")
	}
	fmt.Println()

	fmt.Printf("%-15s", "Worked Days: ")
	fmt.Printf("%-20d", fullStatistics.WorkedDays)
	fmt.Printf("%-15s", "Avg Start: ")
	fmt.Printf("%-20s", fullStatistics.AvgStart)
	fmt.Printf("%-16s", "Sick Days: ")
	fmt.Printf("%-20d\n", fullStatistics.SickDays)

	fmt.Printf("%-15s", "Weekdays: ")
	fmt.Printf("%-20d", fullStatistics.TotalWeekDays)
	fmt.Printf("%-15s", "Avg End: ")
	fmt.Printf("%-20s", fullStatistics.AvgEnd)
	fmt.Printf("%-16s", "Vacation Days: ")
	fmt.Printf("%-20d\n", fullStatistics.VacationDays)

	fmt.Printf("%-15s", "Worked Time: ")
	fmt.Printf("%-20s", fullStatistics.WorkedTime)
	fmt.Printf("%-15s", "Avg Lunch: ")
	fmt.Printf("%-20.2f", fullStatistics.AvgLunch)
	fmt.Printf("%-16s", "Overtime Days: ")
	fmt.Printf("%-20d\n", fullStatistics.OverTimeDays)

	fmt.Printf("%-70s", "")
	fmt.Printf("%-16s", "Total Overtime: ")
	fmt.Printf("%-20.2f\n", fullStatistics.TotalOvertime.Float64)

	fmt.Printf("%-70s", "")
	fmt.Printf("%-16s", "Avg Overtime: ")
	fmt.Printf("%-20.2f\n", fullStatistics.AvgOvertime.Float64)

	fmt.Println()
}

/*func PrintVacSicDaysDiagram(aYear []*data.MonthStats, year int) {
	var maxDays int
	for _, month := range aYear {
		if month.VacationDays > maxDays {
			maxDays = month.VacationDays
		}
		if month.SickDays > maxDays {
			maxDays = month.SickDays
		}
	}

	fmt.Printf("%-5s", "")
	fmt.Print(year, " (Vacation Days ", "██", ", Sickdays ░)", "\n")

	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	fmt.Printf("%-5s", "")
	for _, month := range months {
		fmt.Printf("%-4s", month)
	}
	fmt.Println()

	fmt.Printf("%-4s", "")
	for range 48 {
		fmt.Printf("_")
	}
	fmt.Print("\n")

	for i := maxDays - 1; i >= 0; i-- {
		fmt.Printf("%-3d", i+1)
		fmt.Print("|")
		for j := 0; j < 12; j++ {
			if len(aYear) >= j+1 && aYear[j].VacationDays >= i+1 {
				fmt.Printf("%3s", "██")
			} else {
				fmt.Printf("%3s", "  ")
			}
			if len(aYear) >= j+1 && aYear[j].SickDays >= i+1 {
				fmt.Print("░")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("|\n")
	}
}*/

func PrintWorkedDaysDiagram(monthStats []*data.MonthStats) {
	var maxDays int
	for _, month := range monthStats {
		if month.WorkedDays > maxDays {
			maxDays = month.WorkedDays
		}
		if month.TotalWeekDays > maxDays {
			maxDays = month.TotalWeekDays
		}
	}

	fmt.Printf("%-5s", "")
	fmt.Print("(Worked Days ", "██", ", Weekdays ░)", "\n")

	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	fmt.Printf("%-5s", "")
	for _, month := range months {
		fmt.Printf("%-4s", month)
	}
	fmt.Println()

	fmt.Printf("%-5s", "")
	for idx, _ := range months {
		if len(monthStats) >= idx+1 {
			fmt.Printf("%-4d", monthStats[idx].WorkedDays)
		} else {
			fmt.Printf("%-4d", 0)
		}
	}
	fmt.Println()

	fmt.Printf("%-4s", "")
	for range 48 {
		fmt.Printf("_")
	}
	fmt.Print("\n")

	for i := maxDays - 1; i >= 0; i-- {
		fmt.Printf("%-3d", i+1)
		fmt.Print("|")
		for j := 0; j < 12; j++ {
			if len(monthStats) >= j+1 && monthStats[j].WorkedDays >= i+1 {
				fmt.Printf("%3s", "██")
			} else {
				fmt.Printf("%3s", "  ")
			}
			if len(monthStats) >= j+1 && monthStats[j].TotalWeekDays >= i+1 {
				fmt.Print("░")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("|\n")
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
