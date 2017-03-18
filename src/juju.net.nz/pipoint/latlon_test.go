package pipoint

import (
	"math"
	"testing"
	"fmt"
)

func checkClose(t *testing.T, a, b, limit float64) {
	if math.IsNaN(a) || math.IsNaN(b) {
		t.Errorf("Either %v or %v is NaN", a, b)
	}
	err := math.Abs(a - b)
	if err >= limit {
		t.Errorf("%v differs from %v by %v (more than %v)\n",
			a, b, err, limit)
	}
}

func TestLatLength(t *testing.T) {
	// Test vectors from https://en.wikipedia.org/wiki/Latitude
	checkClose(t, LatLength(0*math.Pi/180), 110574, 1)
	checkClose(t, LatLength(15*math.Pi/180), 110649, 1)
	checkClose(t, LatLength(30*math.Pi/180), 110852, 1)
	checkClose(t, LatLength(45*math.Pi/180), 111132, 1)
	checkClose(t, LatLength(60*math.Pi/180), 111412, 1)
	checkClose(t, LatLength(75*math.Pi/180), 111618, 1)
}

func ExampleLatLength() {
	fmt.Println(int(LatLength(46.8 * math.Pi/180)))
	// Output: 111166
}

func TestLonLength(t *testing.T) {
	// Test vectors from https://en.wikipedia.org/wiki/Latitude
	checkClose(t, LonLength(0*math.Pi/180), 111320, 1)
	checkClose(t, LonLength(15*math.Pi/180), 107550, 1)
	checkClose(t, LonLength(30*math.Pi/180), 96486, 1)
	checkClose(t, LonLength(45*math.Pi/180), 78847, 1)
	checkClose(t, LonLength(60*math.Pi/180), 55800, 1)
	checkClose(t, LonLength(75*math.Pi/180), 28902, 1)
}
