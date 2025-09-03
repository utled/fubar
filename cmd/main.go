package main

import (
	"fTime/cli"
	"fmt"
)

func main() {
	err := cli.Main()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
