package pipoint

import (
	"math"
)

const (
	a = 6378137.0 // in m
	e = 298.257223563
)

// LatLength returns the length of a line of latitude at the given
// latitude.  Input is in rad, output in m.
func LatLength(lat float64) float64 {
	return 111132.954 - 559.822*math.Cos(2*lat) + 1.175*math.Cos(4*lat)
}

// LonLength returns the length of a line of longitude at the given
// latitude.  Input is in rad, output in m.
func LonLength(lat float64) float64 {
	sin := math.Sin(lat)
	return math.Pi * a * math.Cos(lat) /
		(180 * math.Sqrt(1-e*e*sin*sin))
}
