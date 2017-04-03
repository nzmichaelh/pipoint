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

// Geographic coordinates.
type Position struct {
	Time    float64
	Lat     float64
	Lon     float64
	Alt     float64
	Heading float64
}

// On-surface coordinates.
type NEUPosition struct {
	Time  float64
	North float64
	East  float64
	Up    float64
}

// Orientation of a body.
type Attitude struct {
	Roll  float64
	Pitch float64
	Yaw   float64
}

func (p *Position) ToNEU() *NEUPosition {
	lat := AsRad(p.Lat)

	return &NEUPosition{
		Time:  p.Time,
		North: LatLength(lat) * p.Lat,
		East:  LonLength(lat) * p.Lon,
		Up:    p.Alt,
	}
}

func (p *NEUPosition) Sub(right *NEUPosition) *NEUPosition {
	return &NEUPosition{
		Time:  p.Time - right.Time,
		North: p.North - right.North,
		East:  p.East - right.East,
		Up:    p.Up - right.Up,
	}
}

func (p *NEUPosition) Add(right *NEUPosition) *NEUPosition {
	return &NEUPosition{
		Time:  p.Time + right.Time,
		North: p.North + right.North,
		East:  p.East + right.East,
		Up:    p.Up + right.Up,
	}
}
