package main

import (
  "crypto/md5"
	"database/sql"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
  "fmt"
  "strings"
)

type LtarMetObservation struct {
	LTARSiteAcronym       string
	StationId             int
	Datetime              time.Time
	RecordType            string
	Air_temp107_avg       sql.NullFloat64
	Wind_speed_wvt        sql.NullFloat64
	Wind_direction_d1_wvt        sql.NullFloat64
	Relative_humidity_avg sql.NullFloat64
	Rain_mm               sql.NullFloat64
	AirPressure           sql.NullFloat64
	PAR                   sql.NullFloat64
	ShortWaveIn           sql.NullFloat64
	LongWaveIn            sql.NullFloat64
	BatteryVoltage        sql.NullFloat64
	LoggerTemperature     sql.NullFloat64
}

func (d *LtarMetObservation) to_csv() []string {
	values := []string{
		"KBS",
		"000",
		d.Datetime.Format(time.RFC3339Nano),
		"L",
		floatToString(d.Air_temp107_avg),
		floatToString(d.Wind_speed_wvt),
		floatToString(d.Wind_direction_d1_wvt),
		floatToString(d.Relative_humidity_avg),
		floatToString(d.Rain_mm),
		"",
		"",
		"",
		"",
		floatToString(d.BatteryVoltage),
		"",
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
		"WindDirection",
		"RelativeHumidity",
		"Precipitation",
		"AirPressure",
		"PAR",
		"ShortWaveIn",
		"LongWaveIn",
		"BatteryVoltage",
		"LoggerTemperatuure",
	}
	return values
}

func ltar_met_observations(db *sqlx.DB, c *gin.Context) {

	rows, err := db.Queryx(" select * from ( select Air_temp107_avg, Relative_humidity_avg ,Wind_speed_wvt, Wind_direction_d1_wvt, raingauge_hourly.rain_mm,Datetime, Battery_voltage_min as BatteryVoltage from weather.lter_hour_d join weather.raingauge_hourly on raingauge_hourly.hours = lter_hour_d.datetime where datetime < now() - interval '1 hour' order by datetime desc limit $1) t1 order by datetime", limit(c, 1154))

	if err != nil {
		log.Print("error in query")
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
	writer := csv.NewWriter(c.Writer)

	obs := LtarMetObservation{}
	writer.Write(obs.header())
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

func ltar_met_md5(db *sqlx.DB, c *gin.Context) {
	rows, err := db.Queryx(" select * from ( select Air_temp107_avg, Relative_humidity_avg ,Wind_speed_wvt, Wind_direction_d1_wvt, raingauge_hourly.rain_mm,Datetime, Battery_voltage_min as BatteryVoltage from weather.lter_hour_d join weather.raingauge_hourly on raingauge_hourly.hours = lter_hour_d.datetime where datetime < now() - interval '1 hour' order by datetime desc limit $1) t1 order by datetime", limit(c, 1154))

	if err != nil {
		log.Print("error in query")
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
  var data strings.Builder
	writer := csv.NewWriter(&data)

	obs := LtarMetObservation{}
	writer.Write(obs.header())
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
  md5value := md5.Sum([]byte(data.String()))
  md5_slice := fmt.Sprintf("%x\n", md5value[:])

  c.Writer.Write([]byte(md5_slice))
}
