package main

import (
	"database/sql"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type DayObservation struct {
	Year_rtm                     sql.NullInt64
	Day_rtm                      sql.NullInt64
	Hour_rtm                     sql.NullInt64
	Air_temp_107_max             sql.NullFloat64
	Air_temp_107_hour_max        sql.NullInt64
	Air_temp_107_min             sql.NullFloat64
	Air_temp_107_hour_min        sql.NullInt64
	Rh_max                       sql.NullFloat64
	Rh_hour_max                  sql.NullInt64
	Rh_min                       sql.NullFloat64
	Rh_hour_min                  sql.NullInt64
	Solar_radiation              sql.NullFloat64
	Sol_rad_max                  sql.NullFloat64
	Sol_rad_hour_max             sql.NullInt64
	Soil_temp_5_cm_max           sql.NullFloat64
	Soil_temp_5_cm_hour_max      sql.NullInt64
	Soil_temp_5_cm_min           sql.NullFloat64
	Soil_temp_5_cm_hour_min      sql.NullInt64
	Soil_temp_10_cm_max          sql.NullFloat64
	Soil_temp_10_cm_hour_max     sql.NullInt64
	Soil_temp_10_cm_min          sql.NullFloat64
	Soil_temp_10_cm_hour_min     sql.NullInt64
	Soil_moisture_10_cm_max      sql.NullFloat64
	Soil_moisture_10_cm_hour_max sql.NullInt64
	Soil_moisture_10_cm_min      sql.NullFloat64
	Soil_moisture_10_cm_hour_min sql.NullInt64
	Soil_moisture_25_cm_max      sql.NullFloat64
	Soil_moisture_25_cm_hour_max sql.NullInt64
	Soil_moisture_25_cm_min      sql.NullFloat64
	Soil_moisture_25_cm_hour_min sql.NullInt64
	Wind_speed_max               sql.NullFloat64
	Wind_speed_hour_max          sql.NullInt64
	Rain_mm                      sql.NullFloat64
	Battery_voltage_min          sql.NullFloat64
	Date                         time.Time
	Precip                       sql.NullFloat64
}

func (d *DayObservation) toMawn() []string {
	value := []string{
		"24",
		intToString(d.Year_rtm),
		intToString(d.Day_rtm),
		intToString(d.Hour_rtm),
		floatToString(d.Air_temp_107_max),
		intToString(d.Air_temp_107_hour_max),
		floatToString(d.Air_temp_107_min),
		intToString(d.Air_temp_107_hour_min),
		floatToString(d.Rh_max),
		intToString(d.Rh_hour_max),
		floatToString(d.Rh_min),
		intToString(d.Rh_hour_min),
		floatToString(d.Solar_radiation),
		floatToString(d.Sol_rad_max),
		intToString(d.Sol_rad_hour_max),
		floatToString(d.Soil_temp_5_cm_max),
		intToString(d.Soil_temp_5_cm_hour_max),
		floatToString(d.Soil_temp_5_cm_min),
		intToString(d.Soil_temp_5_cm_hour_min),
		floatToString(d.Soil_temp_10_cm_max),
		intToString(d.Soil_temp_10_cm_hour_max),
		floatToString(d.Soil_temp_10_cm_min),
		intToString(d.Soil_temp_10_cm_hour_min),
		floatToString(d.Soil_moisture_10_cm_max),
		intToString(d.Soil_moisture_10_cm_hour_max),
		floatToString(d.Soil_moisture_10_cm_min),
		intToString(d.Soil_moisture_10_cm_hour_min),
		floatToString(d.Soil_moisture_25_cm_max),
		intToString(d.Soil_moisture_25_cm_hour_max),
		floatToString(d.Soil_moisture_25_cm_min),
		intToString(d.Soil_moisture_25_cm_hour_min),
		floatToString(d.Wind_speed_max),
		intToString(d.Wind_speed_hour_max),
		floatToString(d.Rain_mm),
		floatToString(d.Battery_voltage_min),
		d.Date.String(),
	}
	return value
}

func (d *DayObservation) mawnHeader() []string {
	value := []string{
		"#code",
		"year",
		"day",
		"report time",
		"air temperature maximum",
		"time of air temperature maximum",
		"air temperature minimum",
		"time of air temperature minimum",
		"relative humidity maximum",
		"time of relative humidity maximum",
		"relative humidity minimum",
		"time of relative humidity minimum",
		"solar radiation",
		"solar radiation density maximum",
		"time  of solar radiation desnsity",
		"soil temperature maximum at 5 cm under bare soil",
		"time  of soil temperature maximum at 5 cm under bare soil",
		"soil temperature minimum at 5 cm under bare soil",
		"time of soil temperature minimum at 5 cm under bare soil",
		"soil temperature maximum at 10 cm under bare soil",
		"time of soil temperature at 10 cm  under bare soil",
		"soil temperature minimum at 10 cm under bare soil",
		"time of soil temperature minimum at 10 cm under bare soil",
		"soil moisture maximum at 5 cm under sod",
		"time of soil moisture maximum at 5 cm under sod",
		"soil moisture minimum at 5 cm under sod",
		"time of soil moisture minimum at 5 cm under sod",
		"soil moisture maximum at 20 cm under sod",
		"time of soil moisture maximum at 20 cm under sod",
		"soil moisture minimum at  20 cm under sod",
		"time of soil moisture minimum at 20 cm under sod",
		"maximum wind speed",
		"time of maximum wind speed",
		"total precipitation",
		"minimum data logger battery voltage",
		"timestamp",
	}
	return value
}

func (d *DayObservation) mawnUnit() []string {
	value := []string{
		"#",
		"#",
		"", "", "",
		"C", "", "C", "",
		"%", "", "%", "",
		"KJ/m^2", "KJ/m^2/min",
		"", "C", "", "C", "",
		"C", "", "C", "",
		"%", "", "%",
		"", "%", "", "%", "",
		"m/s", "", "mm", "V",
	}
	return value
}

func day_observations(db *sqlx.DB, c *gin.Context) {
	rows, err := db.Queryx("select * from (select * from weather.day_observations_cache order by date desc limit $1) t1 order by date", limit(c))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
	writer := csv.NewWriter(c.Writer)

	observation := DayObservation{}
	writer.Write(observation.mawnHeader())
	writer.Write(observation.mawnUnit())

	for rows.Next() {
		if err := rows.StructScan(&observation); err != nil {
			log.Fatal(err)
		}

		observation.Solar_radiation.Float64 = observation.Solar_radiation.Float64 * 86.4
		observation.Sol_rad_max.Float64 = observation.Sol_rad_max.Float64 * (0.6977 * 60)
		observation.Rh_max.Float64 = observation.Rh_max.Float64 * 100
		observation.Rh_min.Float64 = observation.Rh_min.Float64 * 100

		writer.Write(observation.toMawn())

		if i%500 == 0 {
			writer.Flush()
		}
		i = i + 1

	}
	writer.Flush()
}
