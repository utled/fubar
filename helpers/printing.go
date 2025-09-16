package helpers

import (
	"fTime/data"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func PrintHeader(withSupportText bool, state *data.ReportState) {
	fmt.Print("\n                   ┏━━━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━━┓\n" +
		"                   ┃┏┓┏┓┃┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┓┏┓┃\n" +
		"                   ┗┛┃┃┗┛┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃┃┃\n" +
		"                   ━━┃┃━━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃┃┃\n" +
		"                   ━┏┛┗┓━┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┛┗┛┃\n" +
		"                   ━┗━━┛━━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━━┛\n" +
		"                   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"                   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	if withSupportText {
		fmt.Print("                   cmd -> Display available commands           Ctrl+C -> Close program\n\n")
	}
	if !state.ReportUpToDate {
		fmt.Println("                   There are missing regitrations.\n"+
			"                   Last completed date:", state.MaxCompletedDate)
	}
	fmt.Print("\n                   Selected date:", state.SelectedDate, "\n\n")

}

func PrintSelectedDate(state *data.ReportState) {
	err := ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}
	PrintHeader(true, state)
	fmt.Println(state.SelectedRecord)
}

func PrintCommands(state *data.ReportState) {
	err := ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}
	PrintHeader(false, state)
	availableCommands := []string{
		"\n_____DISPLAY_____________________________________________________________________________________________",
		"today                                     -> Display current date",
		"range <YYYYMMDD YYYYMMDD>             -> Display date range as table",
		"switch <YYYYMMDD>                       -> Switch display date",
		"\n_____DAILY ACTIONS_______________________________________________________________________________________",
		"start <MMSS>                              -> Set/modify start time for displayed date",
		"end <MMSS> [optional]<ot/-ot>>                   -> Set/modify end time for displayed date",
		"ot/-ot                                    -> Set/remove excess time as overtime for displayed date",
		"lunch <INT(minutes)>                      -> Set/modify lunch time for displayed date",
		"addit <INT(minutes)>                      -> Set/modify additional time for displayed date>",
		"off/vac/sic                               -> Flag displayed date as non working day",
		"\n_____SCHEDULING__________________________________________________________________________________________",
		"sched <YYYYMMDD> <YYYYMMDD> <off/vac/sic> -> Schedule date period for coming off time",
		"sched remove                              -> Remove scheduled period",
		"back <norm/off/vac/sic>                        -> Backfill all non registered days back to last completed date",
		"\n_____DEFAULT CONFIGURATIONS______________________________________________________________________________",
		"conflunch <INT(minutes)>                  -> Update default lunch duration",
		"conflength <MMSS>                         -> Update default length of day",
		"\n_____GENERAL_____________________________________________________________________________________________",
		"cmd                                       -> Display available commands",
		"Ctrl+C                                    -> Close program",
		"\n_____GENERAL_____________________________________________________________________________________________",
		"delete                                    -> Delete selected date (can only delete last registered date)",
	}

	for _, command := range availableCommands {
		fmt.Println(command)
	}
	fmt.Println()
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
