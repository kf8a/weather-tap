package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/sqltocsv"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"
)

func CampbellTime(myTime time.Time) (int, int, int) {
	hourmin := myTime.Hour()*100 + myTime.Minute()
	doy := myTime.YearDay()
	if hourmin == 0 {
		hourmin = 2400
		doy = doy - 1
	}
	return myTime.Year(), doy, hourmin
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

type Datum struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

type Variate struct {
	Id      int
	Title   string
	Y_label string
}

func variates(db *sqlx.DB, c *gin.Context) {
	var variates []Variate
	db.Select(&variates, "select id, title, y_label from weather.variates where updated is true")

	obj := gin.H{"variates": variates}
	c.HTML(200, "variates.html", obj)
}

func variatesById(db *sqlx.DB, c *gin.Context) {
	idString := c.Params.ByName("id")
	var query string
	db.Exec("set search_path=weather")
	id, err := strconv.Atoi(idString)
	if err == nil {
		var data = []Datum{}
		err = db.Get(&query, "select query from weather.variates where id = $1", id)
		if err == nil {
			db.Select(&data, query)
		}
		c.JSON(200, data)
	}
}

type Table struct {
	Id    int
	Name  string
	Query string
}

func tables(db *sqlx.DB, c *gin.Context) {
	db.Exec("set search_path=weather")
	var tables []Table
	db.Select(&tables, "select id, name from weather.tables")

	obj := gin.H{"tables": tables}
	c.HTML(200, "tables.html", obj)
}

func tablesById(db *sqlx.DB, c *gin.Context) {
	id := c.Params.ByName("id")
	var query Table
	db.Exec("set search_path=weather")
	err := db.Get(&query, "select id, name, query from weather.tables where id = $1", id)
	if err == nil {
		rows, _ := db.Query(query.Query)
		defer rows.Close()

		w := c.Writer
		w.Header().Set("Content-type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+query.Name+".csv\"")

		sqltocsv.Write(w, rows)
	}
}

func floatToString(value sql.NullFloat64) string {
	result := "nil"
	if value.Valid {
		result = strconv.FormatFloat(value.Float64, 'f', 5, 64)
	}
	return result
}

func intToString(value sql.NullInt64) string {
	result := "nil"
	if value.Valid {
		result = strconv.FormatInt(value.Int64, 10)
	}
	return result
}

func limit(c *gin.Context, limit int) int {
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
	templates := template.Must(template.ParseFiles("templates/tables.html", "templates/variates.html"))
	router.SetHTMLTemplate(templates)

	// router.Static("/assets", "/Users/bohms/code/go/src/weather-tap/assets")
	router.Static("/assets", "./assets")

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
	router.GET("/five_minute_observations.js", func(c *gin.Context) {
		five_minute_observations_js(db, c)
	})
	router.GET("/five_minute_observations.xml", func(c *gin.Context) {
		five_minute_observations_xml(db, c)
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

	connection := "user=" + u.Name + " password=" + u.Password + " dbname=metadata host=granby.kbs.msu.edu port=5432"
	// connection := "user=" + u.Name + " password=" + u.Password + " dbname=metadata host=127.0.0.1 port=5430"
	db, err := sqlx.Open("postgres", connection)
	checkErr(err, "sql.Open failed")
	defer db.Close()

	Router(db).Run(":9000")
}
