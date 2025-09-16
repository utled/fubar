package data

import (
	"database/sql"
)

type OffDay struct {
	OffDate string
	OffType string
}

type ReportState struct {
	ReportUpToDate   bool
	MaxDate          string
	MaxCompletedDate string
	SelectedDate     string
	SelectedRecord   *WorkDateRecord
	ProjectedEnd     string
}

type WorkDateRecord struct {
	WorkDate       string
	StartTime      sql.NullString
	EndTime        sql.NullString
	LunchDuration  sql.NullInt16
	DayTotal       sql.NullString
	DayBalance     sql.NullFloat64
	Overtime       sql.NullBool
	MovingBalance  sql.NullFloat64
	AdditionalTime sql.NullInt16
	DayLength      sql.NullString
	DayType        sql.NullString
}

type UserConfig struct {
	DefaultLunch     sql.NullInt16
	DefaultDayLength sql.NullString
	OffStart         sql.NullString
	OffEnd           sql.NullString
	OffType          sql.NullString
}
