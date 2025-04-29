package main

import (
	"database/sql"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"time"
)

type HourObservation struct {
	Year_rtm                    int
	Day_rtm                     int
	Hourminute_rtm              int
	Air_temp107_avg             sql.NullFloat64
	Relative_humidity_avg       sql.NullFloat64
	Solar_radiation_avg         sql.NullFloat64
	Soil_temperature_5_cm_bare  sql.NullFloat64
	Soil_temperature_10_cm_bare sql.NullFloat64
	Soil_moisture_5_cm          sql.NullFloat64
	Soil_moisture_20_cm         sql.NullFloat64
	Wind_direction_d1_wvt       sql.NullFloat64
	Wind_speed_wvt              sql.NullFloat64
	Rain_mm                     sql.NullFloat64
	Battery_voltage_min         sql.NullFloat64
	Datetime                    time.Time
  Campbell_time               time.Time
}

func (d *HourObservation) toMawn() []string {
	values := []string{
		"60",
		strconv.Itoa(d.Year_rtm),
		strconv.Itoa(d.Day_rtm),
		strconv.Itoa(d.Hourminute_rtm),
		floatToString(d.Air_temp107_avg),
		floatToString(d.Relative_humidity_avg),
		"nil",
		// floatToString(d.Solar_radiation_avg),
		floatToString(d.Soil_temperature_5_cm_bare),
		floatToString(d.Soil_temperature_10_cm_bare),
		floatToString(d.Soil_moisture_5_cm),
		floatToString(d.Soil_moisture_20_cm),
		floatToString(d.Wind_direction_d1_wvt),
		floatToString(d.Wind_speed_wvt),
		"nil", "nil",
		floatToString(d.Rain_mm),
		"nil", "nil",
		floatToString(d.Battery_voltage_min),
		d.Datetime.Format(time.RFC3339),
		d.Campbell_time.Format(time.RFC3339),
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
		"soil moisture at 5 cm",
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
    "campbell_time",
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
    "",
    "",
	}
	return values
}

func hour_observations(db *sqlx.DB, c *gin.Context) {
	// rows, err := db.Queryx("select * from ( select Air_temp107_avg,Relative_humidity_avg,Solar_radiation_avg, Soil_temperature_5_cm_bare, Soil_temperature_10_cm_bare, Soil_moisture_5_cm,Soil_moisture_20_cm,Wind_direction_d1_wvt, Wind_speed_wvt,coalesce(raingauge_hourly.rain_mm, 0) as rain_mm,Battery_voltage_min,Datetime from weather.lter_hour_d left outer join weather.raingauge_hourly on raingauge_hourly.hours = lter_hour_d.datetime where datetime < now() - interval '1 hour' order by datetime desc limit $1) t1 order by datetime", limit(c, 97))

  query := `select * from (SELECT air_temp107_avg, relative_humidity_avg * 100::double precision AS relative_humidity_avg,
  solar_radiation_avg * 0.6977::double precision * 3600::double precision AS solar_radiation_avg,
  soil_temperature_5_cm_bare, soil_temperature_10_cm_bare,
  soil_moisture_5_cm, soil_moisture_20_cm, wind_direction_d1_wvt, wind_speed_wvt,
  rain_tipping_mm as rain_mm, battery_voltage_min,
  datetime ,
  datetime + '1 hour' as campbell_time
  FROM
  weather.lter_hour_d_delete
  WHERE
  datetime < (now() - interval '1 hour')
  and datetime > '2007-12-01'
  order by datetime desc limit $1) t1 order by datetime`

  rows, err := db.Queryx(query, limit(c, 97))

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
	writer := csv.NewWriter(c.Writer)

	obs := HourObservation{}
	writer.Write(obs.mawnHeader())
	writer.Write(obs.mawnUnit())
	for rows.Next() {
		if err := rows.StructScan(&obs); err != nil {
			log.Fatal(err)
		}

		obs.Year_rtm, obs.Day_rtm, obs.Hourminute_rtm = CampbellTime(obs.Campbell_time.Local())

		obs.Solar_radiation_avg.Float64 = obs.Solar_radiation_avg.Float64 * 0.6977 * 3600

		obs.Relative_humidity_avg.Float64 = obs.Relative_humidity_avg.Float64 * 100

		writer.Write(obs.toMawn())

		if i%500 == 0 {
			writer.Flush()
		}
		i = i + 1

	}
	writer.Flush()

}
