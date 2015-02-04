package main

import (
    "github.com/gin-gonic/gin"
    "time"
)

func CampbellTime(myTime time.Time) [3]int {
	hourmin := myTime.Hour()*100 + myTime.Minute()
	if hourmin == 0 {
		hourmin = 2400
	}
	return [3]int{myTime.Year(), myTime.YearDay(), hourmin}
}

func Index(c *gin.Context){
  c.String(200,"we are here")
}

func Hello(c *gin.Context) {
    c.String(200,"hello %s", c.Params.ByName("id"))
}

func VariatesById(c *gin.Context) {
   id := c.Params.ByName("id")
   format := c.Params.ByName("format")
   c.String(200, "variates/%s.%s", id, format)
}

func main() {
	router := gin.Default()
	router.GET("/tables", Index)
	router.GET("/tables/:id", Hello)
	router.GET("/variates", Index)
	router.GET("/variates/:id", VariatesById)
	router.GET("/variates/:id.%format", VariatesById)
	router.GET("/day_observations.:format", Index)
  router.GET("/day_observations/:id.:format", Hello)
	router.GET("/hour_observations.:format", Index)
  router.GET("/hour_observations/:id.:format", Hello)
	router.GET("/five_minute_observations.:format", Index)
  router.GET("/five_minute_observations/:id.:format", Hello)

  router.Run("127.0.0.1:9000")
}
