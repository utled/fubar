package helpers

import (
	"fmt"
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
