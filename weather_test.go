package main

import (
	/* "encoding/base64" */
	"github.com/stretchr/testify/assert"
	"log"
	/* "net/http" */
	/* "net/http/httptest" */
	"testing"
	"time"
)

var timeTest = []struct {
	in  string
	out [3]int
}{
	{
		in:  "2014-01-30T10:30:06+05:00",
		out: [3]int{2014, 30, 1030},
	},
	{
		in:  "2013-02-03T00:00:00+05:00",
		out: [3]int{2013, 34, 2400},
	},
	{
		in:  "2013-05-01T01:54:00+05:00",
		out: [3]int{2013, 121, 154},
	},
}

func TestCambpellTime(t *testing.T) {
	for _, test := range timeTest {
		mytime, err := time.Parse(time.RFC3339, test.in)
		if err != nil {
			log.Fatal(err)
		}
		actual := CampbellTime(mytime)
		assert.Equal(t, test.out, actual)
	}
}

/* GET    /tables(.:format) */
/* GET    /tables/:id(.:format) */
/* GET    /variates(.:format) */
/* GET    /variates/:id(.:format) */
/* GET    /day_observations(.:format) */
/* GET    /day_observations/:id(.:format) */
/* GET    /hour_observations(.:format) */
/* GET    /hour_observations/:id(.:format) */
/* GET    /five_minute_observations(.:format) */
/* GET    /five_minute_observations/:id(.:format) */

/* func TestVariateRoute(t *testing.T) { */
/* 	req, _ := http.NewRequest("GET", "/variates") */
/* } */
