package cli

import (
	"bufio"
	"fTime/home_helpers"
	"fmt"
	"os"
	"strings"
	"time"
)

func checkPreviousCompletion() (isCompleted bool, err error) {
	todayString := time.Now().Format("2006-01-02")
	today, err := time.Parse("2006-01-02", todayString)
	if err != nil {
		return false, fmt.Errorf("failed to parse todays date%v", err)
	}

	maxCompletedString, err := home_helpers.GetMaxCompletedDate()
	if err != nil {
		return false, err
	}
	maxCompleted, err := time.Parse("2006-01-02", maxCompletedString)
	if err != nil {
		return false, fmt.Errorf("failed to parse max completed date%v", err)
	}

	var previousCompleted bool
	dateDiff := today.Sub(maxCompleted)
	if dateDiff.Hours() > 24 {
		previousCompleted = false
	} else {
		previousCompleted = true
	}

	return previousCompleted, nil
}

func home() {
	err := home_helpers.ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}

	//selectedDate := time.Now().Format("2006-01-02")
	selectedDate := "2024-02-10"
	selectedDateRecord, err := home_helpers.GetOneWorkDateRecord(selectedDate)
	if err != nil {
		fmt.Println(err)
	}

	home_helpers.GetTimesheet()

	home_helpers.PrintHeader(true, selectedDate)

	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		arguments := strings.Split(strings.TrimSpace(input), " ")
		switch arguments[0] {
		case "today":
			previousCompleted, err := checkPreviousCompletion()
			if err != nil {
				return
			}
			if previousCompleted {
				fmt.Println("Previous dates are completed correctly")
			} else {
				fmt.Println("Previous dates are not completed correctly")
			}
			fmt.Println(selectedDateRecord)
		case "range":
			fmt.Println("not implemented...")
		case "switch":
			fmt.Println("not implemented...")
		case "start":
			fmt.Println("not implemented...")
		case "end":
			fmt.Println("not implemented...")
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
			err := home_helpers.ClearTerminal()
			if err != nil {
				fmt.Println(err)
			}
			home_helpers.PrintHeader(false, selectedDate)
			home_helpers.PrintCommands()
		default:
			err := home_helpers.ClearTerminal()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Invalid command")
		}
	}
}
