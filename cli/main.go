package cli

import (
	"bufio"
	"fTime/data"
	"fTime/helpers"
	"fTime/registration"
	"fmt"
	"os"
	"strings"
)

func Main() {
	err := helpers.ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}

	userConfig, err := data.GetUserConfig()
	if err != nil {
		fmt.Println(err)
	}

	//selectedDate := time.Now().Format("2006-01-02")
	selectedDate := "2024-12-08"

	currentState, err := setNewState(selectedDate, &userConfig)
	if err != nil {
		fmt.Println(err)
	}

	helpers.PrintHeader(true, &currentState)

	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		arguments := strings.Split(strings.TrimSpace(input), " ")
		switch arguments[0] {
		case "today":
			if len(arguments) == 1 {
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}
		case "range":
			fmt.Println("not implemented...")
		case "switch":
			if len(arguments) == 2 {
				selectedDate, err = helpers.FormatValidDateString(arguments[1])
				if err != nil {
					fmt.Println(err)
					break
				}
				currentState, err = setNewState(selectedDate, &userConfig)
				if err != nil {
					return
				}

				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}
		case "start":
			if len(arguments) == 2 {
				err = registration.RegisterStart(arguments[1], &currentState, &userConfig)
				if err != nil {
					fmt.Println(err)
					break
				}
				currentState, err = setNewState(selectedDate, &userConfig)
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}
		case "end":
			if len(arguments) == 2 {
				err = registration.RegisterEnd(arguments[1], &currentState, &userConfig)
				if err != nil {
					fmt.Println(err)
					break
				}
				currentState, err = setNewState(selectedDate, &userConfig)
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintSelectedDate(&currentState)
			} else if len(arguments) == 3 {
				if arguments[2] == "ot" {
					currentState.SelectedRecord.Overtime.Bool = true
				} else if arguments[2] == "-ot" {
					currentState.SelectedRecord.Overtime.Bool = false
				} else {
					fmt.Println("Invalid argument")
					break
				}
				err = registration.RegisterEnd(arguments[1], &currentState, &userConfig)
				if err != nil {
					fmt.Println(err)
					break
				}
				currentState, err = setNewState(selectedDate, &userConfig)
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}

		case "ot":
			fmt.Println("not implemented...")
		case "-ot":
			fmt.Println("not implemented...")
		case "lunch":
			if len(arguments) == 2 {
				err = registration.RegisterLunch(arguments[1], &currentState)
				if err != nil {
					fmt.Println(err)
					break
				}
				currentState, err = setNewState(selectedDate, &userConfig)
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}

		case "addit":
			if len(arguments) == 2 {
				err = registration.RegisterAdditionalTime(arguments[1], &currentState)
				if err != nil {
					fmt.Println(err)
					break
				}
				currentState, err = setNewState(selectedDate, &userConfig)
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}
		case "off", "vac", "sic":
			if len(arguments) == 1 {
				err = registration.RegisterOffDay(&userConfig, &currentState, arguments[0])
				if err != nil {
					fmt.Println(err)
					break
				}
				currentState, err = setNewState(selectedDate, &userConfig)
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}
		case "sched":
			if len(arguments) == 2 && arguments[1] == "remove" {
				err = registration.RemoveScheduledOffPeriod()
				if err != nil {
					fmt.Println(err)
					break
				}
				userConfig, err = data.GetUserConfig()
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintSelectedDate(&currentState)
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
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}
		case "back":
			if len(arguments) == 1 {
				err = registration.RegisterBackflush(&currentState, arguments[1])
				if err != nil {
					fmt.Println(err)
					break
				}
				currentState, err = setNewState(selectedDate, &userConfig)
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}

		case "conflunch":
			if len(arguments) == 2 {
				err = registration.UpdateDefaultLunch(arguments[1])
				if err != nil {
					fmt.Println(err)
					break
				}
				userConfig, err = data.GetUserConfig()
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}
		case "conflength":
			if len(arguments) == 2 {
				err = registration.UpdateDefaultLength(arguments[1])
				if err != nil {
					fmt.Println(err)
					break
				}
				userConfig, err = data.GetUserConfig()
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}
		case "cmd":
			helpers.PrintCommands(&currentState)
		case "delete":
			if len(arguments) == 1 {
				err = registration.DeleteDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}
		default:
			err := helpers.ClearTerminal()
			if err != nil {
				fmt.Println(err)
			}
			helpers.PrintSelectedDate(&currentState)
			fmt.Println("Invalid command")
		}
	}
}
