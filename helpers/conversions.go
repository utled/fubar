package helpers

import (
	"fTime/utils"
	"fmt"
	"math"
	"strconv"
	"time"
)

func FormatValidDateString(dateString string) (formattedDateString string, err error) {

	if len(dateString) != 8 {
		return "", fmt.Errorf("input must in format <YYYYMMDD>")
	}
	for i := 0; i < len(dateString); i++ {
		if !(dateString[i] >= '0' && dateString[i] <= '9') {
			return "", fmt.Errorf("input includes non-numerical characters")
		}
	}

	year := dateString[0:4]
	month := dateString[4:6]
	day := dateString[6:]
	formattedDateString = year + "-" + month + "-" + day

	return formattedDateString, nil
}

func FormatValidTimeString(timeString string) (formattedTimeString string, err error) {
	if len(timeString) != 4 {
		return "", fmt.Errorf("input must in format <HHMM>")
	}
	hour, err := strconv.Atoi(timeString[0:2])
	if err != nil {
		return "", fmt.Errorf("input contains non-numerical characters")
	}
	if hour >= 24 {
		return "", fmt.Errorf("hour must be less than 24")
	}
	minute, err := strconv.Atoi(timeString[2:4])
	if err != nil {
		return "", fmt.Errorf("input contains non-numerical characters")
	}
	if minute >= 60 {
		return "", fmt.Errorf("minute must be less than 60")
	}

	formattedTimeString = timeString[0:2] + ":" + timeString[2:4] + ":" + "00"

	return formattedTimeString, nil
}

func ParseTimeObject(timeString string) (time.Time, error) {
	dummyDate, err := time.Parse(utils.DateLayout, utils.NonsenseDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse dummy date %v", err)
	}
	registeredTime, err := time.Parse(utils.TimeLayout, timeString)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse registered time %v", err)
	}

	year, month, day := dummyDate.Date()
	hour, minute, second := registeredTime.Clock()
	timeObject := time.Date(year, month, day, hour, minute, second, 0, time.Local)

	return timeObject, nil
}

func DecimalToTime(totalBalance float64) (balanceString string) {
	sign := ""
	if totalBalance < 0 {
		sign = "-"
		totalBalance = math.Abs(totalBalance)
	}
	balanceDuration := time.Duration(totalBalance * float64(time.Hour))
	hours := int(balanceDuration.Hours())
	minutes := int(balanceDuration.Minutes()) % 60
	return fmt.Sprintf("%s%02d:%02d\n", sign, hours, minutes)
}
