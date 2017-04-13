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
	"fmt"
	"log"
	"math"
)

// RunState executes when the camera is tracking the rover.
type RunState struct {
	pi *PiPoint
}

// Update is called when a param is updated.
func (s *RunState) Update(param *Param) {
	switch param {
	case s.pi.neu:
		s.pi.rover.Set(param.Get())
	}

	if !s.pi.rover.Ok() || !s.pi.base.Ok() {
		// Location is invalid or old.
		log.Println("run: skipping as invalid or old", s.pi.rover.Ok(), s.pi.base.Ok())
		return
	}

	rover := s.pi.rover.Get().(*NEUPosition)
	base := s.pi.base.Get().(*NEUPosition)
	offset := s.pi.baseOffset.Get().(*NEUPosition)

	att, err := point(rover, base, offset)
	if err != nil {
		log.Printf("point: %v\n", err)
		return
	}

	if param == s.pi.neu {
		offset := s.pi.offset.Get().(*Attitude)
		s.pi.pan.Set(WrapAngle(att.Yaw + offset.Yaw))
		s.pi.tilt.Set(WrapAngle(att.Pitch + offset.Pitch))
	}
}

func point(rover, base, offset *NEUPosition) (*Attitude, error) {
	delta := rover.Sub(base.Add(offset))
	if math.Abs(delta.North) > 10e3 || math.Abs(delta.East) > 10e3 {
		return nil, fmt.Errorf("Rover is too far away")
	}

	hdist := math.Sqrt(delta.North*delta.North + delta.East*delta.East)
	pitch := math.Atan2(delta.Up, hdist)
	yaw := math.Atan2(delta.East, delta.North)

	return &Attitude{
		Pitch: pitch,
		Yaw:   yaw,
	}, nil
}
