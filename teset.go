package main

import (
	"fmt"
	"time"
)

func main() {
	currentDate := time.Now()
	doesItWork := currentDate.Weekday() == time.Sunday
	fmt.Println("Who's work?", doesItWork)
}
