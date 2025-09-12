package actions

import (
	"fTime/helpers"
	"fmt"
	"strconv"
)

func UpdateDefaultLunch(lunchDurationString string) error {
	lunchDurationInt, err := strconv.Atoi(lunchDurationString)
	if err != nil {
		return fmt.Errorf("failed to convert input to numeric value.\nInput format must be <MM>")
	}
	lunchDuration := int16(lunchDurationInt)
	if lunchDuration < 0 {
		return fmt.Errorf("lunchDuration can't be a negative value")
	}

	err = helpers.UpdateDefaultLunch(lunchDuration)
	if err != nil {
		return err
	}

	return nil
}

func UpdateDefaultLength(dayLengthString string) error {
	dayLengthInt, err := strconv.Atoi(dayLengthString)
	if err != nil {
		return fmt.Errorf("failed to convert input to numeric value.\nInput format must be <MM>")
	}
	dayLength := int16(dayLengthInt)
	if dayLength < 0 {
		return fmt.Errorf("lunchDuration can't be a negative value")
	}

	err = helpers.UpdateDefaultLength(dayLength)
	if err != nil {
		return err
	}

	return nil
}
