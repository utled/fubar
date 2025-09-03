package home_helpers

import (
	"fTime/db"
	"fmt"
	"github.com/go-gota/gota/dataframe"
)

func GetTimesheet() {
	timesheetData, err := db.GetTimesheetData()
	if err != nil {
		fmt.Println(err)
	}

	var timesheet []map[string]interface{}
	columns, _ := timesheetData.Columns()

	for timesheetData.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for val := range values {
			pointers[val] = &values[val]
		}

		if err := timesheetData.Scan(pointers...); err != nil {
			fmt.Println("failed to scan values", err)
		}

		row := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]
			if val != nil {
				row[colName] = val
			}
		}
		timesheet = append(timesheet, row)

	}
	timesheetDF := dataframe.LoadMaps(timesheet)
	timesheetDF = timesheetDF.Arrange(dataframe.RevSort("workdate"))
	fmt.Println(timesheetDF)
}
