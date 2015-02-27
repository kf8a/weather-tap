package main

import (
	"database/sql"
	"time"
)

type HourObservation struct {
	Year_rtm sql.NullInt64
	Datetime time.Time
}

func (d *HourObservation) mawnHeader() []string {
	values := []string{
		"#code",
		"year",
		"day of year",
		"report hour minute",
		"air temperature",
		"relative humidity",
		"solar radiation",
		"soil temperature at 5 cm",
		"soil temperature at 10 cm",
		"soil  moisture at 5 cm",
		"soil moisture at 20 cm",
		"wind direction",
		"wind speed",
		"maximum hourly wind speed",
		"time of maximum wind speed",
		"precipitation",
		"leaf0",
		"leaf1",
		"battery voltage minimum",
		"timestamp",
	}
	return values
}

func (d *HourObservation) mawnUnit() []string {
	values := []string{
		"#",
		"",
		"",
		"",
		"C",
		"%",
		"kJ/m^2",
		"C",
		"C",
		"%",
		"%",
		"degrees",
		"m/s",
		"m/s",
		"",
		"mm",
		"",
		"",
		"",
		"V",
	}
	return values
}
