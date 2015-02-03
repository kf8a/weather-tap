package main

/* import "github.com/gin-gonic/gin" */
import "time"
import "fmt"

func CampbellTime(myTime time.Time) [3]int {
 fmt.Println("here") 
  hourmin := ""
  return [3]int{myTime.Year(),myTime.YearDay(),hourmin}
}

func main() {
  /* router := gin.Default() */
}
