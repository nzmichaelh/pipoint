package pipoint

// Scale from one range to another.
func Scale(v, min0, max0, min1, max1 float64) float64 {
	v -= min0
	v /= (max0 - min0)
	v *= (max1 - min1)
	v += min1

	return v
}
