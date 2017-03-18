package main

import (
	"math"
)

const (
	a = 6378137.0 // m
	e = 298.257223563
)

func LatLength(lat float64) float64 {
	return 111132.954-559.822*math.Cos(2*lat) + 1.175*math.Cos(4*lat)
}

func LonLength(lat float64) float64 {
	sin := math.Sin(lat)
	return math.Pi * a * math.Cos(lat) /
		(180*math.Sqrt(1-e*e*sin*sin))
}
