package main

import (
	"database/sql"
	"strconv"
	"time"
)

type HourObservation struct {
	Year_rtm              int
	Day_rtm               int
	Hourminute_rtm        int
	Air_temp107_avg       sql.NullFloat64
	Relative_humidity_avg sql.NullFloat64
	Solar_radiation_avg   sql.NullFloat64
	Soil_temp_q_avg       sql.NullFloat64
	Soil_moisture_5_cm    sql.NullFloat64
	Soil_moisture_20_cm   sql.NullFloat64
	Wind_direction_d1_wvt sql.NullFloat64
	Wind_speed_wvt        sql.NullFloat64
	Rain_mm               sql.NullFloat64
	Battery_voltage_min   sql.NullFloat64
	Datetime              time.Time
}

func (d *HourObservation) toMawn() []string {
	values := []string{
		"60",
		strconv.Itoa(d.Year_rtm),
		strconv.Itoa(d.Day_rtm),
		strconv.Itoa(d.Hourminute_rtm),
		floatToString(d.Air_temp107_avg),
		floatToString(d.Relative_humidity_avg),
		floatToString(d.Solar_radiation_avg),
		floatToString(d.Soil_temp_q_avg),
		floatToString(d.Soil_moisture_5_cm),
		floatToString(d.Soil_moisture_20_cm),
		floatToString(d.Wind_direction_d1_wvt),
		floatToString(d.Wind_speed_wvt),
		floatToString(d.Rain_mm),
		floatToString(d.Battery_voltage_min),
		d.Datetime.Format(time.RFC3339),
	}
	return values
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
