package pipoint

import (
	"math"
)

// Scale from one range to another.
func Scale(v, min0, max0, min1, max1 float64) float64 {
	v -= min0
	v /= (max0 - min0)
	v *= (max1 - min1)
	v += min1

	return v
}

func WrapAngle(v float64) float64 {
	// TODO: add tests.
	for ; v > math.Pi; {
		v -= math.Pi*2
	}
	for ; v < -math.Pi; {
		v += math.Pi*2
	}
	return v
}