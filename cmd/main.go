package main

import (
	"fTime/cli"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println()
		os.Exit(0)
	}()

	err := cli.Main()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

}
