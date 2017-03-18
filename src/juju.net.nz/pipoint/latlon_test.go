package main

import (
	"math"
	"testing"
)

func checkClose(t *testing.T, a float64, b float64, limit float64) {
	err := math.Abs(a - b)
	if err >= limit {
		t.Errorf("%v differs from %v by %v (more than %v)\n",
			a, b, err, limit)
	}
}

func TestLatLength(t *testing.T) {
	// Test vectors from https://en.wikipedia.org/wiki/Latitude
	checkClose(t, LatLength(math.Pi*0/180), 110574, 1)
	checkClose(t, LatLength(math.Pi*15/180), 110649, 1)
	checkClose(t, LatLength(math.Pi*30/180), 110852, 1)
	checkClose(t, LatLength(math.Pi*45/180), 111132, 1)
	checkClose(t, LatLength(math.Pi*60/180), 111412, 1)
	checkClose(t, LatLength(math.Pi*75/180), 111618, 1)
}

func TestLonLength(t *testing.T) {
	// Test vectors from https://en.wikipedia.org/wiki/Latitude
	checkClose(t, LonLength(math.Pi*0/180), 111320, 1)
	checkClose(t, LonLength(math.Pi*15/180), 107550, 1)
	checkClose(t, LonLength(math.Pi*30/180), 96486, 1)
	checkClose(t, LonLength(math.Pi*45/180), 78847, 1)
	checkClose(t, LonLength(math.Pi*60/180), 55800, 1)
	checkClose(t, LonLength(math.Pi*75/180), 28902, 1)
}
