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
	exitChan := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	err := cli.Main()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	go func() {
		<-signalChan
		close(exitChan)
	}()
}
