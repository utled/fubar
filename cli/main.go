package cli

import (
	"bufio"
	"fTime/actions"
	"fTime/helpers"
	"fmt"
	"os"
	"strings"
)

func Main() {
	err := helpers.ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}

	userConfig, err := helpers.GetUserConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("User Config:", userConfig)

	//selectedDate := time.Now().Format("2006-01-02")
	selectedDate := "2024-12-08"

	currentState, err := setNewState(selectedDate)
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
				currentState, err = setNewState(selectedDate)
				if err != nil {
					return
				}

				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}
		case "start":
			if len(arguments) == 2 {
				err = actions.RegisterStart(arguments[1], &currentState)
				if err != nil {
					fmt.Println(err)
					break
				}
				currentState, err = setNewState(selectedDate)
				if err != nil {
					fmt.Println(err)
				}
				helpers.PrintSelectedDate(&currentState)
			} else {
				fmt.Println("Invalid argument")
			}
		case "end":
			/*if len(arguments) == 2 {
				err = actions.RegisterEnd(arguments[1], false, &currentState)
				if err != nil {
					fmt.Println(err)
					break
				}
			} else if len(arguments) == 3 && arguments[2] == "ot" {
				err = actions.RegisterEnd(arguments[1], true, &currentState)
				if err != nil {
					fmt.Println(err)
					break
				}
			} else {
				fmt.Println("Invalid argument")
			}*/

		case "ot":
			fmt.Println("not implemented...")
		case "-ot":
			fmt.Println("not implemented...")
		case "lunch":
			fmt.Println("not implemented...")
		case "addit":
			fmt.Println("not implemented...")
		case "off":
			fmt.Println("not implemented...")
		case "vac":
			fmt.Println("not implemented...")
		case "sic":
			fmt.Println("not implemented...")
		case "sched":
			fmt.Println("not implemented...")
		case "back":
			fmt.Println("not implemented...")
		case "conflunch":
			fmt.Println("not implemented...")
		case "conflength":
			fmt.Println("not implemented...")
		case "cmd":
			helpers.PrintCommands(&currentState)
		case "test":

		default:
			err := helpers.ClearTerminal()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Invalid command")
		}
	}
}
