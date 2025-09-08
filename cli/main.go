package cli

import (
	"bufio"
	"fTime/helpers"
	"fmt"
	"os"
	"strings"
)

type reportStatus struct {
	reportUpToDate   bool
	maxCompletedDate string
	selectedDate     string
	selectedRecord   helpers.WorkDateRecord
}

func (ws *reportStatus) GetReportUpToDate() bool {
	return ws.reportUpToDate
}

func (ws *reportStatus) GetMaxCompletedDate() string {
	return ws.maxCompletedDate
}

func (ws *reportStatus) GetSelectedDate() string {
	return ws.selectedDate
}

func (ws *reportStatus) GetSelectedRecord() helpers.WorkDateRecord {
	return ws.selectedRecord
}

func Main() {
	err := helpers.ClearTerminal()
	if err != nil {
		fmt.Println(err)
	}

	//selectedDate := time.Now().Format("2006-01-02")
	selectedDate := "2024-12-08"

	recordExists, err := helpers.CheckIfDateExists(selectedDate)
	if err != nil {
		fmt.Println(err)
	}

	var selectedDateRecord helpers.WorkDateRecord
	if recordExists {
		selectedDateRecord, err = helpers.GetOneWorkDateRecord(selectedDate)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		selectedDateRecord = helpers.WorkDateRecord{
			WorkDate: selectedDate,
		}
	}

	previousCompleted, maxCompletedDate, err := helpers.CheckPreviousCompletion()
	if err != nil {
		return
	}

	workingStatus := reportStatus{
		reportUpToDate:   previousCompleted,
		maxCompletedDate: maxCompletedDate,
		selectedDate:     selectedDate,
		selectedRecord:   selectedDateRecord,
	}

	helpers.PrintHeader(true, &workingStatus)

	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		arguments := strings.Split(strings.TrimSpace(input), " ")
		switch arguments[0] {
		case "today":
			if len(arguments) == 1 {
				helpers.PrintSelectedDate(&workingStatus)
			} else {
				fmt.Println("Invalid argument")
			}
		case "range":
			fmt.Println("not implemented...")
		case "switch":
			if len(arguments) == 2 {
				selectedDate, err := helpers.FormatValidDateString(arguments[1])
				if err != nil {
					fmt.Println(err)
					break
				}
				workingStatus.selectedDate = selectedDate
				selectedDateRecord, err = helpers.GetOneWorkDateRecord(selectedDate)
				if err != nil {
					fmt.Println(err)
				}
				workingStatus.selectedRecord = selectedDateRecord

				helpers.PrintSelectedDate(&workingStatus)
			} else {
				fmt.Println("Invalid argument")
			}
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
			helpers.PrintCommands(&workingStatus)
		default:
			err := helpers.ClearTerminal()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Invalid command")
		}
	}
}
