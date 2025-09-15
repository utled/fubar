package registration

import (
	"fTime/data"
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

	err = data.UpdateDefaultLunch(lunchDuration)
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

	err = data.UpdateDefaultLength(dayLength)
	if err != nil {
		return err
	}

	return nil
}

func ScheduleOffPeriod(
	offStart string,
	offEnd string,
	offType string,
	config *data.UserConfig,
) error {
	if config.OffStart.String != "" {
		return fmt.Errorf("scheduled off period already exists.\nuse <sched remove> to remove existing period")
	}

	if offType != "off" && offType != "vac" && offType != "sic" {
		return fmt.Errorf("off period needs to be either 'off', 'vac' or 'sic'")
	}

	offStartFormatted, err := helpers.FormatValidDateString(offStart)
	if err != nil {
		return err
	}

	offEndFormatted, err := helpers.FormatValidDateString(offEnd)
	if err != nil {
		return err
	}

	dateIsInFuture, err := helpers.CheckDateInFuture(offStartFormatted)
	if err != nil {
		return err
	}
	if !dateIsInFuture {
		return fmt.Errorf("off period must be in the future")
	}

	startIsBeforeEnd, err := helpers.CheckDateBefore(offStartFormatted, offEndFormatted)
	if err != nil {
		return err
	}
	if !startIsBeforeEnd {
		return fmt.Errorf("start date must be before end date")
	}

	err = data.UpdateScheduledOff(offStartFormatted, offEndFormatted, offType)
	if err != nil {
		return err
	}

	return nil
}

func RemoveScheduledOffPeriod() error {
	err := data.UpdateScheduledOff("", "", "")
	if err != nil {
		return err
	}

	return nil
}
