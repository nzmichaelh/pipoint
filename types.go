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

// Position is a 3D point in geographic coordinates.
type Position struct {
	Time    float64
	Lat     float64
	Lon     float64
	Alt     float64
	Heading float64
}

// NEUPosition is a 3D point on the local tangent plane.
type NEUPosition struct {
	Time  float64
	North float64
	East  float64
	Up    float64
}

// Attitude is the orientation of a body, often to the local tangent
// plane.
type Attitude struct {
	Roll  float64
	Pitch float64
	Yaw   float64
}

// ToNEU converts a geographic position to local tangent plane.
func (p *Position) ToNEU() *NEUPosition {
	lat := AsRad(p.Lat)

	return &NEUPosition{
		Time:  p.Time,
		North: LatLength(lat) * p.Lat,
		East:  LonLength(lat) * p.Lon,
		Up:    p.Alt,
	}
}

// Sub returns piecewise this minus right.
func (p *NEUPosition) Sub(right *NEUPosition) *NEUPosition {
	return &NEUPosition{
		Time:  p.Time - right.Time,
		North: p.North - right.North,
		East:  p.East - right.East,
		Up:    p.Up - right.Up,
	}
}

// Add returns piecewise this plus right.
func (p *NEUPosition) Add(right *NEUPosition) *NEUPosition {
	return &NEUPosition{
		Time:  p.Time + right.Time,
		North: p.North + right.North,
		East:  p.East + right.East,
		Up:    p.Up + right.Up,
	}
}
