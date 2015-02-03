package main

import (
	"github.com/stretchr/testify/assert"
	"log"
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
