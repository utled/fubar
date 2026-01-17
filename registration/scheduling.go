package registration

import (
	"fmt"
	"fubar/data"
	"fubar/helpers"
)

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

	startIsWknd, err := helpers.CheckIfDateIsWknd(offStartFormatted)
	if err != nil {
		return err
	}
	endIsWknd, err := helpers.CheckIfDateIsWknd(offEndFormatted)
	if err != nil {
		return err
	}
	if startIsWknd || endIsWknd {
		return fmt.Errorf("off period can't start or end on a wknd")
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
