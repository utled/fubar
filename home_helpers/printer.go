package home_helpers

import "fmt"

func PrintHeader(withSupportText bool, selectedDate string) {
	fmt.Print("\n                   ┏━━━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━┓━━━┓\n" +
		"                   ┃┏┓┏┓┃┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┫┣┛┓┏┓┃\n" +
		"                   ┗┛┃┃┗┛┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃┃┃\n" +
		"                   ━━┃┃━━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃━┃┃┃┃\n" +
		"                   ━┏┛┗┓━┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┫┣┓┛┗┛┃\n" +
		"                   ━┗━━┛━━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━┛━━━┛\n" +
		"                   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"                   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	/*	fmt.Print("\n ▄▄▄▄▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄ ▄▄▄▄▄▄  \n" +
		"█       █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █      █ \n" +
		"█▄     ▄█   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █  ▄    █\n" +
		"  █   █ █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █ █ █   █\n" +
		"  █   █ █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █ █▄█   █\n" +
		"  █   █ █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █   █       █\n" +
		"  █▄▄▄█ █▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄█▄▄▄▄▄▄█ \n\n")*/

	if withSupportText {
		fmt.Print("                   cmd -> Display available commands           Ctrl+C -> Close program\n\n")
	}
	fmt.Print("\n                   Date:", selectedDate, "\n\n")

}

func PrintSelectedDate() {
	fmt.Println("not implemented...")
}

func PrintCommands() {
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
