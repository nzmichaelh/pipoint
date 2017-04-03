// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
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
