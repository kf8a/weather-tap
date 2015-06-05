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

type LtarMetObservation struct {
	LTARSiteAcronym       string
	StationId             int
	Datetime              time.Time
	RecordType            string
	Air_temp107_avg       sql.NullFloat64
	Wind_speed_wvt        sql.NullFloat64
	Relative_humidity_avg sql.NullFloat64
	Rain_mm               sql.NullFloat64
}

func (d *LtarMetObservation) to_csv() []string {
	values := []string{
		"KBS",
		"000",
		d.Datetime.Format(time.RFC3339Nano),
		"L",
		floatToString(d.Air_temp107_avg),
		floatToString(d.Wind_speed_wvt),
		floatToString(d.Relative_humidity_avg),
		floatToString(d.Rain_mm),
	}
	return values
}

func (d *LtarMetObservation) header() []string {
	values := []string{
		"#LTARSiteAcronym",
		"StationId",
		"DateTime",
		"RecordType",
		"AirTemperature",
		"WindSpeed",
		"RelativeHumidity",
		"Precipitation",
	}
	return values
}

func (d *LtarMetObservation) units() []string {
	values := []string{
		"#",
		"",
		"",
		"",
		"C",
		"m/s",
		"%",
		"mm",
		"kPa",
	}
	return values
}

func ltar_met_observations(db *sqlx.DB, c *gin.Context) {

	rows, err := db.Queryx(" select * from ( select Air_temp107_avg,Relative_humidity_avg ,Wind_speed_wvt, raingauge_hourly.rain_mm,Datetime from weather.lter_hour_d join weather.raingauge_hourly on raingauge_hourly.hours = lter_hour_d.datetime where datetime < now() - interval '1 hour' order by datetime desc limit $1) t1 order by datetime", limit(c, 1154))

	if err != nil {
		log.Print("error in query")
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
	writer := csv.NewWriter(c.Writer)

	obs := LtarMetObservation{}
	writer.Write(obs.header())
	writer.Write(obs.units())
	for rows.Next() {
		if err := rows.StructScan(&obs); err != nil {
			log.Fatal(err)
		}

		obs.Relative_humidity_avg.Float64 = obs.Relative_humidity_avg.Float64 * 100

		writer.Write(obs.to_csv())

		if i%500 == 0 {
			writer.Flush()
		}
		i = i + 1

	}
	writer.Flush()
}
