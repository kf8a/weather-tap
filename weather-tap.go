package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

func CampbellTime(myTime time.Time) [3]int {
	hourmin := myTime.Hour()*100 + myTime.Minute()
	if hourmin == 0 {
		hourmin = 2400
	}
	return [3]int{myTime.Year(), myTime.YearDay(), hourmin}
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

type DayObservationDB struct {
	Battery_voltage_min          sql.NullFloat64
	Year_rtm                     sql.NullInt64
	Day_rtm                      sql.NullInt64
	Hour_rtm                     sql.NullInt64
	Air_temp_107_max             sql.NullFloat64
	Air_temp_107_hour_max        sql.NullInt64
	Air_temp_107_min             sql.NullFloat64
	Air_temp_107_hour_min        sql.NullInt64
	Rh_max                       sql.NullFloat64 // rh_max * 100  float
	Rh_hour_max                  sql.NullInt64
	Rh_min                       sql.NullFloat64 //rh_min * 100 float
	Rh_hour_min                  sql.NullInt64
	Sol_rad_max                  sql.NullFloat64 //sol_rad_max * (0.6977 * 60)
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
	Date                         time.Time
	Solar_radiation              sql.NullFloat64
	Precip                       sql.NullFloat64
}

func day_observations(db *sqlx.DB, c *gin.Context) {
	rows, err := db.Queryx("select * from weather.day_observations_cache")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		observation := DayObservationDB{}
		if err := rows.StructScan(&observation); err != nil {
			log.Fatal(err)
		}

		observation.Sol_rad_max.Float64 = observation.Sol_rad_max.Float64 * (0.6977 * 60)
		observation.Rh_max.Float64 = observation.Rh_max.Float64 * 100
		observation.Rh_min.Float64 = observation.Rh_min.Float64 * 100
		data, _ := json.Marshal(observation)
		c.Writer.Write(data)

	}
	/* solar_radiation = nil */
	/* if solar_radiation */
	/*   solar_radiation = self.solar_radiation * 86.4 */
	/* end */

	/* sprintf("24,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s, %s, %s,%s, %s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s, %s,%s,%s,%s, %s,%s, %s, %s, %s", */
	/*         year_rtm, day_rtm,hour_rtm, */
	/*         self.air_temp_107_max, */
	/*         self.air_temp_107_hour_max,  self.air_temp_107_min, */
	/*         self.air_temp_107_hour_min, */
	/*         self.rh_max * 100, self.rh_hour_max, self.rh_min * 100, */
	/*         self.rh_hour_min, solar_radiation, */
	/*         self.sol_rad_max * (0.6977 * 60) , self.sol_rad_hour_max, */
	/*         self.soil_temp_5_cm_max,soil_temp_5_cm_hour_max, */
	/*         soil_temp_5_cm_min, soil_temp_5_cm_hour_min, */
	/*         self.soil_temp_10_cm_max, */
	/*         soil_temp_10_cm_hour_max, */
	/*         self.soil_temp_10_cm_min, */
	/*         soil_temp_10_cm_hour_min, */
	/*         self.soil_moisture_10_cm_max, self.soil_moisture_10_cm_hour_max, */
	/*         self.soil_moisture_10_cm_min, */
	/*         self.soil_moisture_10_cm_hour_min, self.soil_moisture_25_cm_max, */
	/*         self.soil_moisture_25_cm_hour_max, */
	/*         self.soil_moisture_25_cm_min, self.soil_moisture_25_cm_hour_min, */
	/*         self.wind_speed_max, */
	/*         self.wind_speed_hour_max, self.rain_mm, self.battery_voltage_min, */
	/*         self.date) */
}

func hour_observations(db *sqlx.DB, c *gin.Context) {
}

func five_minute_observations(db *sqlx.DB, c *gin.Context) {
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
