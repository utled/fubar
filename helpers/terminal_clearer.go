package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

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
