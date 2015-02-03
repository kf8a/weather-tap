package main

/* import "github.com/gin-gonic/gin" */
import "time"

func CampbellTime(myTime time.Time) [3]int {
	hourmin := myTime.Hour()*100 + myTime.Minute()
	if hourmin == 0 {
		hourmin = 2400
	}
	return [3]int{myTime.Year(), myTime.YearDay(), hourmin}
}

func main() {
	/* router := gin.Default() */
}
