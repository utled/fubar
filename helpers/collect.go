package helpers

import (
	"fmt"
	"fubar/data"
	"strconv"
)

func CollectDailyInputs(
	startTime,
	lunchDuration,
	endTime,
	additionalTime,
	dayType,
	overtime string) (inputValues data.CollectedValues, err error) {
	if startTime != "" {
		inputValues.StartTime, err = FormatValidTimeString(startTime)
		if err != nil {
			return inputValues, fmt.Errorf("invalid start time: must be in format <HH:MM> or <HHMM>")
		}
	}
	if lunchDuration != "" {
		inputValues.LunchDuration, err = strconv.Atoi(lunchDuration)
		if err != nil {
			return inputValues, fmt.Errorf("lunch duration contains non-numerical characters")
		}
	}
	if endTime != "" {
		inputValues.EndTime, err = FormatValidTimeString(endTime)
		if err != nil {
			return inputValues, fmt.Errorf("invalid end time: must be in format <HH:MM> or <HHMM>")
		}
	}
	if additionalTime != "" {
		inputValues.AdditionalTime, err = strconv.Atoi(additionalTime)
		if err != nil {
			return inputValues, fmt.Errorf("additional time contains non-numerical characters")
		}
	}
	inputValues.DayType = dayType
	inputValues.Overtime, err = strconv.ParseBool(overtime)
	if err != nil {
		return inputValues, err
	}

	return inputValues, nil
}
