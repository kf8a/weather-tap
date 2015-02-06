package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
  "os"
  "log"
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

func Hello(c *gin.Context) {
	c.String(200, "hello %s", c.Params.ByName("id"))
}

type Datum struct {
  Time time.Time `json:"date"`
  Value float64   `json:"value"`
}

func jsonify(data []Datum) string {
	json_data, err := json.Marshal(data)
	checkErr(err, "Failed to convert to Json")

  return string(json_data)
}

func VariatesById(db *sqlx.DB, c *gin.Context) {
	id := c.Params.ByName("id")
  var query string
  db.Exec("set search_path=weather")
  db.Get(&query, "select query from weather.variates where id = $1", id)
  var data = []Datum{}
  db.Select(&data, query)
  c.String(200, jsonify(data))
}

func Router(db *sqlx.DB) *gin.Engine {
	router := gin.Default()
	router.GET("/tables", Index)
	router.GET("/tables/:id", Hello)
	router.GET("/variates", Index)
	router.GET("/variates/:id", func(c *gin.Context) {
    VariatesById(db, c)
  })
	router.GET("/day_observations", Index)
	router.GET("/day_observations/:id", func(c *gin.Context) {
    VariatesById(db, c)
  })
	router.GET("/hour_observations", Index)
	router.GET("/hour_observations/:id", Hello)
	router.GET("/five_minute_observations", Index)
	router.GET("/five_minute_observations/:id", Hello)

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

  connection := "user="+u.Name +" password=" + u.Password +" dbname=metadata host=granby.kbs.msu.edu"
	db, err := sqlx.Open("postgres", connection )
	checkErr(err, "sql.Open failed")
	defer db.Close()


	Router(db).Run("127.0.0.1:9000")
}
