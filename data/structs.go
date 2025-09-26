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
	DayType        sql.NullString
	StartTime      sql.NullString
	EndTime        sql.NullString
	LunchDuration  sql.NullInt16
	AdditionalTime sql.NullInt16
	Overtime       sql.NullBool
	DayTotal       sql.NullString
	DayBalance     sql.NullFloat64
	TotalBalance   sql.NullFloat64
	DayLength      sql.NullString
}

type UserConfig struct {
	ID               int
	DefaultLunch     sql.NullInt16
	DefaultDayLength sql.NullString
	OffStart         sql.NullString
	OffEnd           sql.NullString
	OffType          sql.NullString
}

type MonthStats struct {
	Month         string
	TotalWeekDays int
	WorkedDays    int
	WorkedTime    string
	VacationDays  int
	SickDays      int
	WeekendDays   int
	OffDays       int
	OverTimeDays  int
	TotalOvertime sql.NullFloat64
}

type FullStats struct {
	WorkedDays    int
	TotalWeekDays int
	WorkedTime    string
	AvgStart      string
	AvgEnd        string
	AvgLunch      float32
	SickDays      int
	VacationDays  int
	OverTimeDays  int
	TotalOvertime sql.NullFloat64
	AvgOvertime   sql.NullFloat64
}
