package helpers

import (
	"fmt"
	"fubar/utils"
	"math"
	"strconv"
	"strings"
	"time"
)

func FormatValidDateString(dateString string) (formattedDateString string, err error) {
	dateString = strings.Replace(dateString, "-", "", -1)
	dateString = strings.TrimSpace(dateString)
	if len(dateString) != 8 {
		return "", fmt.Errorf("input must in format <YYYYMMDD> or <YYYY-MM-DD>")
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
	timeStringReplaced := strings.Replace(timeString, ":", "", -1)
	timeStringTrimmed := strings.TrimSpace(timeStringReplaced)
	if len(timeStringTrimmed) != 4 {
		return "", fmt.Errorf("input must in format <HH:MM> or <HHMM>")
	}
	hour, err := strconv.Atoi(timeStringTrimmed[0:2])
	if err != nil {
		return "", fmt.Errorf("input contains non-numerical characters")
	}
	if hour >= 24 {
		return "", fmt.Errorf("hour must be less than 24")
	}
	minute, err := strconv.Atoi(timeStringTrimmed[2:4])
	if err != nil {
		return "", fmt.Errorf("input contains non-numerical characters")
	}
	if minute >= 60 {
		return "", fmt.Errorf("minute must be less than 60")
	}

	formattedTimeString = timeStringTrimmed[0:2] + ":" + timeStringTrimmed[2:4] + ":" + "00"

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
	return fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
}

func NextDateFromString(dateString string) (nextDateString string, err error) {
	dateParts := strings.Split(dateString, "-")
	year, err := strconv.Atoi(dateParts[0])
	if err != nil {
		return "", fmt.Errorf("failed to parse year %v", err)
	}
	month, err := strconv.Atoi(dateParts[1])
	if err != nil {
		return "", fmt.Errorf("failed to parse month %v", err)
	}
	day, err := strconv.Atoi(dateParts[2])
	if err != nil {
		return "", fmt.Errorf("failed to parse day %v", err)
	}
	nextDateString = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1).Format(utils.DateLayout)
	return nextDateString, nil
}

func PreviousDateFromString(dateString string) (previousDateString string, err error) {
	dateParts := strings.Split(dateString, "-")
	year, err := strconv.Atoi(dateParts[0])
	if err != nil {
		return "", fmt.Errorf("failed to parse year %v", err)
	}
	month, err := strconv.Atoi(dateParts[1])
	if err != nil {
		return "", fmt.Errorf("failed to parse month %v", err)
	}
	day, err := strconv.Atoi(dateParts[2])
	if err != nil {
		return "", fmt.Errorf("failed to parse day %v", err)
	}
	previousDateString = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1).Format(utils.DateLayout)
	return previousDateString, nil
}
