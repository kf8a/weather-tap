package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/sqltocsv"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// "github.com/prometheus/client_golang/prometheus"
	"html/template"
	"log"
	"math"
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
  Host     string `json:"host"`
  Port     string `json:"port"`
  Database string `json:"database"`
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

func removeNans(data []Datum) []Datum {
	d := make([]Datum, 0)

	for _, v := range data {
		if !math.IsNaN(v.Value) {
			d = append(d, v)
		}
	}
	return d
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
		// Loop over the data and replace NaN's with nil
		data = removeNans(data)
		c.JSON(200, data)
	}
}

type Table struct {
	Id    int
	Name  string
	Query string
  Download bool
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
	err := db.Get(&query, "select id, name, query, download from weather.tables where id = $1", id)
	if err == nil {
		rows, _ := db.Query(query.Query)
		defer rows.Close()

		w := c.Writer
    if query.Download {
      w.Header().Set("Content-type", "text/csv")
      w.Header().Set("Content-Disposition", "attachment; filename=\""+query.Name+".csv\"")
    }

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
  config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://oshtemo.kbs.msu.edu"}
  config.AllowMethods = []string{"GET"}

  router.Use(cors.New(config))
  router.SetTrustedProxies([]string{"127.0.0.1","oshtemo.kbs.msu.edu"})

	templates := template.Must(template.ParseFiles("templates/tables.html", "templates/variates.html"))
	router.SetHTMLTemplate(templates)

	router.Static("/assets", "./assets")
	router.GET("/metrics", func(c *gin.Context) {
    promhttp.Handler()
		// prometheus.Handler()
	})
	router.Static("/weather/assets", "./assets")

	router.GET("/tables", func(c *gin.Context) {
		tables(db, c)
	})
	router.GET("/weather/tables", func(c *gin.Context) {
		tables(db, c)
	})
	router.GET("/tables/:id", func(c *gin.Context) {
		tablesById(db, c)
	})
	router.GET("/weather/tables/:id", func(c *gin.Context) {
		tablesById(db, c)
	})
	router.GET("/variates", func(c *gin.Context) {
		variates(db, c)
	})
	router.GET("/weather/variates", func(c *gin.Context) {
		variates(db, c)
	})
	router.GET("/variates/:id", func(c *gin.Context) {
		variatesById(db, c)
	})
	router.GET("/weather/variates/:id", func(c *gin.Context) {
		variatesById(db, c)
	})
	router.GET("/day_observations.mawn", func(c *gin.Context) {
		day_observations(db, c)
	})
	router.GET("/weather/day_observations.mawn", func(c *gin.Context) {
		day_observations(db, c)
	})
	router.GET("/hour_observations.mawn", func(c *gin.Context) {
		hour_observations(db, c)
	})
	router.GET("/weather/hour_observations.mawn", func(c *gin.Context) {
		hour_observations(db, c)
	})
	router.GET("/five_minute_observations.mawn", func(c *gin.Context) {
		five_minute_observations(db, c)
	})
	router.GET("/weather/five_minute_observations.mawn", func(c *gin.Context) {
		five_minute_observations(db, c)
	})
	router.GET("/five_minute_observations.js", func(c *gin.Context) {
		five_minute_observations_js(db, c)
	})
	router.GET("/weather/five_minute_observations.js", func(c *gin.Context) {
		five_minute_observations_js(db, c)
	})
	router.GET("/five_minute_observations.xml", func(c *gin.Context) {
		five_minute_observations_xml(db, c)
	})
	router.GET("/weather/five_minute_observations.xml", func(c *gin.Context) {
		five_minute_observations_xml(db, c)
	})

	return router
}

func main() {
	configFile, err := os.Open(".env")
	checkErr(err, "can't read .env file")

	u := User{}
	jsonParser := json.NewDecoder(configFile)
  println(jsonParser)
	if err = jsonParser.Decode(&u); err != nil {
		checkErr(err, "parsing config file")
	}

  connection := "user=" + u.Name + " password=" + u.Password + " dbname=" + u.Database +  " host=" + u.Host + " port=" + u.Port
	db, err := sqlx.Open("postgres", connection)
	checkErr(err, "sql.Open failed")
	defer db.Close()

	Router(db).Run(":9000")
}
