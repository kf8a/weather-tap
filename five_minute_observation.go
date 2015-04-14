package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"strconv"
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

type Rain struct {
	Rain_mm  float64   `xml:"rain-mm"`
	Datetime time.Time `xml:"datetime"`
}

func (d *FiveMinuteObservation) toMawn() []string {
	values := []string{
		"5",
		strconv.Itoa(d.Year_rtm),
		strconv.Itoa(d.Day_rtm),
		strconv.Itoa(d.Hourminute_rtm),
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

func five_minute_observations(db *sqlx.DB, c *gin.Context) {

	rows, err := db.Queryx("select * from (select air_temp107_avg, relative_humidity_avg, leaf_wetness_mv_avg, solar_radiation_avg, wind_direction_d1_wvt, wind_speed_wvt, reportPCP * 25.4 as rain_mm, lter_five_minute_a.datetime from weather.lter_five_minute_a left join weather.raingauge on raingauge.datetime = lter_five_minute_a.datetime order by datetime desc limit $1 ) t1 order by datetime", limit(c, 1154))

	if err != nil {
		log.Print("error in query")
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
	writer := csv.NewWriter(c.Writer)

	obs := FiveMinuteObservation{}
	writer.Write(obs.mawnHeader())
	writer.Write(obs.mawnUnit())
	for rows.Next() {
		if err := rows.StructScan(&obs); err != nil {
			log.Fatal(err)
		}

		obs.Year_rtm, obs.Day_rtm, obs.Hourminute_rtm = CampbellTime(obs.Datetime.Local())

		obs.Relative_humidity_avg.Float64 = obs.Relative_humidity_avg.Float64 * 100

		writer.Write(obs.toMawn())

		if i%500 == 0 {
			writer.Flush()
		}
		i = i + 1

	}
	writer.Flush()
}

func five_minute_observations_js(db *sqlx.DB, c *gin.Context) {
	data := []FiveMinuteObservation{}

	db.Select(&data, "select rain_mm, datetime from weather.lter_five_minute_a order by datetime desc limit $1", limit(c, 3))
	c.JSON(200, data)
}

func five_minute_observations_xml(db *sqlx.DB, c *gin.Context) {
	data := []FiveMinuteObservation{}

	db.Select(&data, "select rain_mm, datetime from weather.lter_five_minute_a order by datetime desc limit $1", limit(c, 3))
	output := make([]Rain, len(data))
	for key, value := range data {
		output[key].Rain_mm = value.Rain_mm.Float64
		output[key].Datetime = value.Datetime
	}
	xmlOut, err := xml.MarshalIndent(output, " ", " ")
	if err != nil {
		log.Fatal(err)
	}
	c.Writer.Write(xmlOut)
}
