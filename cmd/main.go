package main

import (
	"fmt"
	"fubar"
	"fubar/cli"
	"fubar/helpers"
	"fubar/tui"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		err := helpers.ClearTerminal()
		if err != nil {
			return
		}
		os.Exit(0)
	}()

	arguments := os.Args
	switch len(arguments) {
	case 1:
		tui.Launch()
	case 2:
		switch arguments[1] {
		case "cli":
			helpers.InitClearFunctions()
			cli.Launch()
		case "test":
			fubar.Test()
		default:
			fmt.Println("invalid argument")
		}
	default:
		fmt.Println("invalid argument")
	}
}
