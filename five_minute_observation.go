package main

import (
	"database/sql"
	"time"
)

type FiveMinuteObservation struct {
	Year_rtm              int
	Day_rtm               int
	Hourminute_rtm        int
	Air_temp107_avg       sql.NullFloat64
	Relative_humidity_avg sql.NullFloat64
	Leaf_wetness_mv_avg   sql.NullFloat64
	Solar_radiation_avg   sql.NullFloat64
	Wind_direction_d1_wvt sql.NullFloat64
	Wind_speed_wvt        sql.NullFloat64
	Rain_mm               sql.NullFloat64
	Datetime              time.Time
}

func (d *FiveMinuteObservation) toMawn() []string {
	values := []string{
		"5",
		string(d.Year_rtm),
		string(d.Day_rtm),
		string(d.Hourminute_rtm),
		floatToString(d.Rain_mm),
		floatToString(d.Leaf_wetness_mv_avg),
		"",
		floatToString(d.Wind_speed_wvt),
		floatToString(d.Air_temp107_avg),
		floatToString(d.Relative_humidity_avg),
		d.Datetime.Format(time.RFC3339),
	}
	return values
}

func (d *FiveMinuteObservation) mawnHeader() []string {
	values := []string{
		"#code",
		"year",
		"day",
		"time",
		"rain_mm",
		"leaf wetness A",
		"leaf wetnetss B",
		"wind speed",
		"air temperature",
		"relative humidity",
		"timestamp",
	}
	return values
}

func (d *FiveMinuteObservation) mawnUnit() []string {
	values := []string{
		"#",
		"",
		"",
		"",
		"mm",
		"",
		"",
		"m/s",
		"C",
		"%",
	}
	return values
}
