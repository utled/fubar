package helpers

import (
	"fTime/utils"
	"fmt"
	"time"
)

func CalcDayTotal(dateRecord *WorkDateRecord) (string, error) {
	startTime, err := ParseTimeObject(dateRecord.StartTime.String)
	if err != nil {
		return "", fmt.Errorf("failed to parse start time")
	}

	endTime, err := ParseTimeObject(dateRecord.EndTime.String)
	if err != nil {
		return "", fmt.Errorf("failed to parse end time")
	}

	timeDiff := endTime.Sub(startTime)
	dayTotal := timeDiff -
		(time.Duration(dateRecord.LunchDuration.Int16) * time.Minute) +
		(time.Duration(dateRecord.AdditionalTime.Int16) * time.Minute)

	totalSeconds := int(dayTotal.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	dayTotalString := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	return dayTotalString, nil
}

func CalcDayBalance(dateRecord *WorkDateRecord) (float64, error) {
	totalTime, err := ParseTimeObject(dateRecord.DayTotal.String)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse day total time")
	}
	fmt.Println(totalTime)

	defaultDayLength, err := ParseTimeObject(dateRecord.DayLength.String)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse day default length")
	}

	timeDiff := totalTime.Sub(defaultDayLength)

	return timeDiff.Hours(), nil
}

func CalcTotalBalance(dateRecord *WorkDateRecord, previousTotal float64) float64 {
	if dateRecord.Overtime.Bool {
		return previousTotal
	}

	newTotalBalance := dateRecord.DayBalance.Float64 + previousTotal

	return newTotalBalance
}

func CalcProjectedEnd(dateRecord *WorkDateRecord, userConfig *UserConfig) string {
	startTime, err := ParseTimeObject(dateRecord.StartTime.String)
	if err != nil {
		return ""
	}

	dayLength, err := ParseTimeObject(dateRecord.DayLength.String)
	if err != nil {
		return ""
	}

	var lunchDuration time.Time
	if dateRecord.LunchDuration.Valid {
		lunchDuration = dayLength.Add(time.Minute * time.Duration(dateRecord.LunchDuration.Int16))
		fmt.Println(lunchDuration)
	} else {
		lunchDuration = dayLength.Add(time.Minute * time.Duration(userConfig.DefaultLunch.Int16))
	}

	addHour := time.Duration(lunchDuration.Hour()) * time.Hour
	addMinute := time.Duration(lunchDuration.Minute()) * time.Minute
	projectedEnd := startTime.Add(addHour + addMinute)

	return projectedEnd.Format(utils.TimeLayout)
}
