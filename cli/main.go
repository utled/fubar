package cli

import (
	"bufio"
	"fmt"
	"fubar/data"
	"fubar/helpers"
	"fubar/registration"
	"fubar/utils"
	"os"
	"strings"
	"time"
)

func Main() {
	err := helpers.ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}

	selectedDate := time.Now().Format(utils.DateLayout)
	currentState := data.ReportState{}
	userConfig, err := data.GetUserConfig()
	if err != nil {
		fmt.Println(err)
	}

	setNewState(selectedDate, &currentState, &userConfig)

	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		arguments := strings.Split(strings.TrimSpace(input), " ")
		switch arguments[0] {
		case "today", "t":
			if len(arguments) == 1 {
				selectedDate = time.Now().Format(utils.DateLayout)
				setNewState(selectedDate, &currentState, &userConfig)
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'today'")
			}
		case "next", "n":
			if len(arguments) == 1 {
				currentDate, err := time.Parse(utils.DateLayout, selectedDate)
				if err != nil {
					fmt.Println(err)
				}
				selectedDate = currentDate.AddDate(0, 0, 1).Format(utils.DateLayout)
				setNewState(selectedDate, &currentState, &userConfig)
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'today'")
			}
		case "prev", "p":
			if len(arguments) == 1 {
				currentDate, err := time.Parse(utils.DateLayout, selectedDate)
				selectedDate = currentDate.AddDate(0, 0, -1).Format(utils.DateLayout)
				if err != nil {
					fmt.Println(err)
				}
				setNewState(selectedDate, &currentState, &userConfig)
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'today'")
			}
		case "switch", "sw":
			if len(arguments) == 2 {
				selectedDate, err = helpers.FormatValidDateString(arguments[1])
				if err != nil {
					fmt.Println(err)
					break
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'switch <YYYYMMDD>'")
			}
		case "last", "la":
			if len(arguments) == 2 {
				timesheet, err := helpers.SetDateRangeFromDayCount(arguments[1])
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintDateRange(timesheet, false, &currentState)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'last <INT(days)>'")
			}
		case "range", "ra":
			if len(arguments) == 3 {
				timesheet, err := helpers.SetDateRangeFromDates(arguments[1], arguments[2])
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintDateRange(timesheet, true, &currentState)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'range <YYYYMMDD YYYYMMDD>'")
			}
		case "clear", "c":
			err = helpers.ClearTerminal()
			if err != nil {
				fmt.Println(err)
			}
			helpers.PrintHeader(&currentState)
		case "start", "s":
			if len(arguments) == 2 {
				err = registration.RegisterStart(arguments[1], &currentState, &userConfig)
				if err != nil {
					fmt.Println(err)
					break
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'start <YYYYMMDD>'")
			}
		case "end", "e":
			if len(arguments) == 2 {
				currentState.SelectedRecord.DayType.String = "norm"
				err = registration.RegisterEnd(arguments[1], &currentState, &userConfig)
				if err != nil {
					fmt.Println(err)
					break
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else if len(arguments) == 3 {
				if arguments[2] == "ot" {
					currentState.SelectedRecord.Overtime.Bool = true
				} else if arguments[2] == "-ot" {
					currentState.SelectedRecord.Overtime.Bool = false
				} else {
					fmt.Println("Invalid argument.\nExpects: 'end <MMSS [optional]ot/-ot'")
					break
				}
				currentState.SelectedRecord.DayType.String = "norm"
				err = registration.RegisterEnd(arguments[1], &currentState, &userConfig)
				if err != nil {
					fmt.Println(err)
					break
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'end <MMSS> [optional]<ot/-ot>'")
			}

		case "ot", "-ot":
			if len(arguments) == 1 {
				err = registration.RegisterOvertime(arguments[0], &currentState, &userConfig)
				if err != nil {
					fmt.Println(err)
					break
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'of/-ot' only")
			}
		case "lunch", "l":
			if len(arguments) == 2 {
				err = registration.RegisterLunch(arguments[1], &currentState)
				if err != nil {
					fmt.Println(err)
					break
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println("Invalid argument.\nExpects: lunch '<INT(minutes)>'")
			}

		case "addit", "ad":
			if len(arguments) == 2 {
				err = registration.RegisterAdditionalTime(arguments[1], &currentState)
				if err != nil {
					fmt.Println(err)
					break
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'addit <INT(minutes)>'")
			}
		case "off", "vac", "sic":
			if len(arguments) == 1 {
				err = registration.RegisterOffDay(&userConfig, &currentState, arguments[0])
				if err != nil {
					fmt.Println(err)
					break
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'off/vac/sic' only")
			}
		case "norm":
			if len(arguments) == 1 {
				err = registration.RevertOffDay(&userConfig, &currentState)
				if err != nil {
					fmt.Println(err)
					break
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'norm' only")
			}
		case "sched", "sc":
			if len(arguments) == 2 {
				switch arguments[1] {
				case "remove":
					err = registration.RemoveScheduledOffPeriod()
					if err != nil {
						fmt.Println(err)
						break
					}
					userConfig, err = data.GetUserConfig()
					if err != nil {
						fmt.Println(err)
					}
					setNewState(selectedDate, &currentState, &userConfig)
				case "show":
					helpers.PrintScheduledOffPeriod(&userConfig, &currentState)
				}
			} else if len(arguments) == 4 {
				err = registration.ScheduleOffPeriod(arguments[1], arguments[2], arguments[3], &userConfig)
				if err != nil {
					fmt.Println(err)
					break
				}
				userConfig, err = data.GetUserConfig()
				if err != nil {
					fmt.Println(err)
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println("Invalid argument.\n" +
					"Expects: 'sched <YYYYMMDD> <YYYYMMDD> <off/vac/sic>' or\n" +
					"'sched remove' or\n" +
					"'sched show'")
			}
		case "back":
			if len(arguments) == 2 {
				err = registration.RegisterBackflush(arguments[1], &currentState, &userConfig)
				if err != nil {
					fmt.Println(err)
					break
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'back <norm/off/vac/sic>'")
			}
		case "conf":
			if len(arguments) == 2 && (arguments[1] == "show" || arguments[1] == "s") {
				helpers.PrintUserConfig(&userConfig, &currentState)
			} else if len(arguments) == 3 && (arguments[1] == "lunch" || arguments[1] == "l") {
				err = registration.UpdateDefaultLunch(arguments[2])
				if err != nil {
					fmt.Println(err)
					break
				}
				userConfig, err = data.GetUserConfig()
				if err != nil {
					fmt.Println(err)
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else if len(arguments) == 3 && (arguments[1] == "length" || arguments[1] == "le") {
				err = registration.UpdateDefaultLength(arguments[2])
				if err != nil {
					fmt.Println(err)
					break
				}
				userConfig, err = data.GetUserConfig()
				if err != nil {
					fmt.Println(err)
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println(
					"Invalid argument.\n" +
						"Expects: " +
						"'conf lunch[l] <INT(minutes)>' or\n" +
						"'conf length[le] <MMSS>' or\n" +
						"'conf show[s]' only")
			}
		case "cmd":
			helpers.PrintCommands(&currentState)
		case "delete", "dl":
			if len(arguments) == 1 {
				err = registration.DeleteDate(&currentState)
				if err != nil {
					fmt.Println(err)
					break
				}
				setNewState(selectedDate, &currentState, &userConfig)
			} else {
				fmt.Println("Invalid argument.\nExpects: 'delete' only")
			}
		case "stats", "st":
			if len(arguments) >= 2 {
				switch arguments[1] {
				case "all", "a":
					if len(arguments) == 2 {
						err = helpers.DisplayAllStatistics(&currentState)
						if err != nil {
							fmt.Println(err)
						}
					} else {
						fmt.Println("Invalid argument.\nExpects: 'stats[st] all[a]' only")
					}
				case "sum", "s":
					if len(arguments) == 3 {
						err = helpers.DisplaySumStatistics(arguments[2], &currentState)
						if err != nil {
							fmt.Println(err)
						}
					} else {
						fmt.Println("Invalid argument.\nExpects: 'stats[st] sum[s] <YYYY>'")
					}
				case "year", "y":
					if len(arguments) == 3 {
						err = helpers.DisplayYearStatistics(arguments[2], &currentState)
						if err != nil {
							fmt.Println(err)
						}
					} else {
						fmt.Println("Invalid argument.\nExpects: 'stats[st] year[y] <YYYY>'")
					}
				case "month", "m":
					if len(arguments) == 4 {
						err = helpers.DisplayMonthStatistics(arguments[2], arguments[3], &currentState)
						if err != nil {
							fmt.Println(err)
						}
					} else {
						fmt.Println("Invalid argument.\nExpects: 'stats[st] month[m] <INT(monthnum)> <YYYY>'")
					}
				case "day", "d":
					if len(arguments) == 3 {
						err = helpers.DisplayDaysStatistics(arguments[2], &currentState)
						if err != nil {
							fmt.Println(err)
						}
					} else {
						fmt.Println("Invalid argument.\nExpects: 'stats[st] day[d] <INT(days)>'")
					}
				case "range", "r":
					if len(arguments) == 4 {
						err = helpers.DisplayRangeStatistics(arguments[2], arguments[3], &currentState)
						if err != nil {
							fmt.Println(err)
						}
					} else {
						fmt.Println("Invalid argument.\nExpects: 'stats[st] range[r] <YYYYMMDD> <YYYYMMDD>'")
					}
				default:
					fmt.Println("Invalid argument.\nRun 'cmd' to display available commands.")
				}
			} else {
				fmt.Println("Invalid argument.\nRun 'cmd' to display available commands.")
			}
		default:
			err := helpers.ClearTerminal()
			if err != nil {
				fmt.Println(err)
			}
			helpers.PrintSelectedDate(&currentState)
			fmt.Println("Invalid command.\nRun 'cmd' to display available commands.")
		}
	}
}
