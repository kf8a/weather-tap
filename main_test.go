package main

import (
	"math"
	"testing"
	"time"
)

func TestNanRemoval(t *testing.T) {
	data := make([]Datum, 2)
	data[0].Time = time.Now()
	data[0].Value = math.NaN()
	data[1].Time = time.Now()
	data[1].Value = 34

	data = removeNans(data)
	if len(data) > 1 {
		t.Error("Expected 1 row got ", len(data))
	}
}
