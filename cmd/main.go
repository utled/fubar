package main

import (
	"fTime/cli"
	"fTime/db"
	"fTime/helpers"
	"fmt"
	"log"
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

	err := db.InitializeDB()
	if err != nil {
		log.Fatalf("error initializing database: %v", err)
	}

	helpers.InitClearFunctions()

	cli.Main()

}
