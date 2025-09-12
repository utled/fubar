package helpers

import (
	"fmt"
)

func PrintHeader(withSupportText bool, state *ReportState) {
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

func PrintSelectedDate(state *ReportState) {
	err := ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}
	PrintHeader(true, state)
	fmt.Println(state.SelectedRecord)
}

func PrintCommands(state *ReportState) {
	err := ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}
	PrintHeader(false, state)
	availableCommands := []string{
		"\n_____DISPLAY_____________________________________________________________________________________________",
		"today                                     -> Display current date",
		"range <YYYY-MM-DD YYYY-MM-DD>             -> Display date range as table",
		"switch <YYYY-MM-DD>                       -> Switch display date",
		"\n_____DAILY ACTIONS_______________________________________________________________________________________",
		"start <MMSS>                              -> Set/modify start time for displayed date",
		"end <MMSS [optional]ot>                   -> Set/modify end time for displayed date",
		"ot/-ot                                    -> Set/remove excess time as overtime for displayed date",
		"lunch <INT(minutes)>                      -> Set/modify lunch time for displayed date",
		"addit <INT(minutes)>                      -> Set/modify additional time for displayed date>",
		"off/vac/sic                               -> Flag displayed date as non working day",
		"\n_____SCHEDULING__________________________________________________________________________________________",
		"sched <YYYY-MM-DD YYYY-MM-DD off/vac/sic> -> Schedule date period for coming off time",
		"back <off/vac/sic>                        -> Backfill all non registered days back to last completed date",
		"\n_____DEFAULT CONFIGURATIONS______________________________________________________________________________",
		"conflunch <INT(minutes)>                  -> Update default lunch duration",
		"conflength <MMSS>                         -> Update default length of day",
		"\n_____GENERAL_____________________________________________________________________________________________",
		"cmd                                       -> Display available commands",
		"Ctrl+C                                    -> Close program",
	}

	for _, command := range availableCommands {
		fmt.Println(command)
	}
	fmt.Println()
}
