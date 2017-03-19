package pipoint

import (
	"math"
)

const (
	m1 = 111132.92 // latitude calculation term 1
	m2 = -559.82   // latitude calculation term 2
	m3 = 1.175     // latitude calculation term 3
	m4 = -0.0023   // latitude calculation term 4
	p1 = 111412.84 // longitude calculation term 1
	p2 = -93.5     // longitude calculation term 2
	p3 = 0.118     // longitude calculation term 3
)

// LatLength returns the length of a line of latitude at the given
// latitude.  Input is in rad, output in m.
func LatLength(lat float64) float64 {
	return m1 + m2*math.Cos(2*lat) + m3*math.Cos(4*lat)
}

// LonLength returns the length of a line of longitude at the given
// latitude.  Input is in rad, output in m.
func LonLength(lat float64) float64 {
	return p1*math.Cos(lat) + p2*math.Cos(3*lat) + p3*math.Cos(5*lat)
}

func AsRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func AsDeg(rad float64) float64 {
	return rad * (180 / math.Pi)
}
