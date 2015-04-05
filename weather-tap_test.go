package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var timeTest = []struct {
	in         string
	year_rtm   int
	day_rtm    int
	hourminute int
}{
	{
		in:         "2014-01-30T10:30:06+05:00",
		year_rtm:   2014,
		day_rtm:    30,
		hourminute: 1030,
	},
	{
		in:         "2013-02-03T00:00:00+05:00",
		year_rtm:   2013,
		day_rtm:    33,
		hourminute: 2400,
	},
	{
		in:         "2013-05-01T01:54:00+05:00",
		year_rtm:   2013,
		day_rtm:    121,
		hourminute: 154,
	},
}

func TestCambpellTime(t *testing.T) {
	for _, test := range timeTest {
		mytime, err := time.Parse(time.RFC3339, test.in)
		if err != nil {
			log.Fatal(err)
		}
		year_rtm, day_rtm, hourminute := CampbellTime(mytime)
		assert.Equal(t, test.year_rtm, year_rtm)
		assert.Equal(t, test.day_rtm, day_rtm)
		assert.Equal(t, test.hourminute, hourminute)
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

func TestVariateRoute(t *testing.T) {
	w := httptest.NewRecorder()
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	r, _ := http.NewRequest("GET", "/variates/1", nil)
	Router(db).ServeHTTP(w, r)
	assert.Equal(t, w.Code, http.StatusOK)
}

// func TestTablesRoute(t *testing.T) {
// 	r, _ := http.NewRequest("GET", "/tables/1", nil)
// 	w := httptest.NewRecorder()

// 	var db *sqlx.DB
// 	Router(db).ServeHTTP(w, r)
// 	assert.Equal(t, w.Code, http.StatusOK)
// }
