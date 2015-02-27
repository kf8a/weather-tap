package main

import (
	"database/sql"
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

func limit(c *gin.Context) int {
	limit := 1500
	query := c.Request.URL.Query()
	if query["limit"] != nil {
		value, err := strconv.Atoi(query["limit"][0])
		if err != nil {
			log.Fatal(err)
		}
		limit = value
	}
	return limit
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
	router.GET("/day_observations.mawn", func(c *gin.Context) {
		day_observations(db, c)
	})
	router.GET("/hour_observations.mawn", func(c *gin.Context) {
		hour_observations(db, c)
	})
	router.GET("/five_minute_observations.mawn", func(c *gin.Context) {
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
