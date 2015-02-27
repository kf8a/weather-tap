package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

func CampbellTime(myTime time.Time) (int, int, int) {
	hourmin := myTime.Hour()*100 + myTime.Minute()
	if hourmin == 0 {
		hourmin = 2400
	}
	return myTime.Year(), myTime.YearDay(), hourmin
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func Index(c *gin.Context) {
	c.String(200, "we are here")
}

type Datum struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

func variatesById(db *sqlx.DB, c *gin.Context) {
	id := c.Params.ByName("id")
	var query string
	db.Exec("set search_path=weather")
	db.Get(&query, "select query from weather.variates where id = $1", id)
	var data = []Datum{}
	db.Select(&data, query)
	c.JSON(200, data)
}

type Table struct {
	Id   int
	Name string
}

type Variate struct {
	Id    int
	Title string
}

func tables(db *sqlx.DB, c *gin.Context) {
	db.Exec("set search_path=weather")
	var tables []Table
	db.Select(&tables, "select id, name from weather.tables")
	c.JSON(200, tables)
}

func variates(db *sqlx.DB, c *gin.Context) {
	db.Exec("set search_path=weather")
	var variates []Variate
	db.Select(&variates, "select id, title from weather.variates")
	c.JSON(200, variates)
}

func tablesById(db *sqlx.DB, c *gin.Context) {
	id := c.Params.ByName("id")
	var query string
	db.Exec("set search_path=weather")
	db.Get(&query, "select query from weather.tables where id = $1", id)
	rows, _ := db.Queryx(query)
	var results []interface{}
	for rows.Next() {
		cols, _ := rows.SliceScan()
		results = append(results, cols)
	}
	rows.Close()
	c.JSON(200, results)
}

func floatToString(value sql.NullFloat64) string {
	return strconv.FormatFloat(value.Float64, 'f', 5, 64)
}

func intToString(value sql.NullInt64) string {
	return strconv.FormatInt(value.Int64, 10)
}

func day_observations(db *sqlx.DB, c *gin.Context) {
	rows, err := db.Queryx("select * from weather.day_observations_cache")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
	writer := csv.NewWriter(c.Writer)

	observation := DayObservation{}
	writer.Write(observation.mawnHeader())
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

func hour_observations(db *sqlx.DB, c *gin.Context) {
	rows, err := db.Queryx("select Air_temp107_avg,Relative_humidity_avg,Solar_radiation_avg,Soil_temp_q_avg,Soil_moisture_5_cm,Soil_moisture_20_cm,Wind_direction_d1_wvt,Wind_speed_wvt,Rain_mm,Battery_voltage_min,Datetime from weather.lter_hour_d")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
	writer := csv.NewWriter(c.Writer)

	obs := HourObservation{}
	writer.Write(obs.mawnHeader())
	for rows.Next() {
		if err := rows.StructScan(&obs); err != nil {
			log.Fatal(err)
		}

		obs.Year_rtm, obs.Day_rtm, obs.Hourminute_rtm = CampbellTime(obs.Datetime.Local())

		obs.Relative_humidity_avg.Float64 = obs.Relative_humidity_avg.Float64 * 100
		obs.Solar_radiation_avg.Float64 = obs.Solar_radiation_avg.Float64 * 0.6977 * 3600

		writer.Write(obs.toMawn())

		if i%500 == 0 {
			writer.Flush()
		}
		i = i + 1

	}
	writer.Flush()

}

func five_minute_observations(db *sqlx.DB, c *gin.Context) {
	rows, err := db.Queryx("select * from weather.lter_five_minute_a")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	i := 0
	writer := csv.NewWriter(c.Writer)

	obs := HourObservation{}
	writer.Write(obs.mawnHeader())
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

func Router(db *sqlx.DB) *gin.Engine {
	router := gin.Default()
	router.GET("/tables", func(c *gin.Context) {
		tables(db, c)
	})
	router.GET("/tables/:id", func(c *gin.Context) {
		tablesById(db, c)
	})
	router.GET("/variates", func(c *gin.Context) {
		variates(db, c)
	})
	router.GET("/variates/:id", func(c *gin.Context) {
		variatesById(db, c)
	})
	router.GET("/day_observations", func(c *gin.Context) {
		day_observations(db, c)
	})
	router.GET("/hour_observations", func(c *gin.Context) {
		hour_observations(db, c)
	})
	router.GET("/five_minute_observations", func(c *gin.Context) {
		five_minute_observations(db, c)
	})

	return router
}

func main() {
	configFile, err := os.Open(".env")
	checkErr(err, "can't read .env file")

	u := User{}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&u); err != nil {
		checkErr(err, "parsing config file")
	}

	connection := "user=" + u.Name + " password=" + u.Password + " dbname=metadata host=127.0.0.1 port=5430"
	db, err := sqlx.Open("postgres", connection)
	checkErr(err, "sql.Open failed")
	defer db.Close()

	Router(db).Run("127.0.0.1:9000")
}
