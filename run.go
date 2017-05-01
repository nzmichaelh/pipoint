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

func (s *RunState) Name() string {
	return "Run"
}

// Update is called when a param is updated.
func (s *RunState) Update(param *Param) {
	switch param {
	case s.pi.neu:
		s.pi.rover.Set(param.Get())
	case s.pi.seconds:
		if (param.GetInt()%5) == 0 && s.pi.vel.Ok() {
			kph := s.pi.vel.GetFloat64() * 3.6
			if kph >= 2 {
				s.pi.audio.Say(fmt.Sprintf("%.0f kph", kph))
			}
		}
	case s.pi.mark:
		s.pi.state.Inc()
	}

	if param != s.pi.rover {
		return
	}

	if !s.pi.rover.Ok() || !s.pi.base.Ok() {
		return
	}

	rover := s.pi.rover.Get().(*NEUPosition)
	base := s.pi.base.Get().(*NEUPosition)
	baseOffset := s.pi.baseOffset.Get().(*NEUPosition)

	att, err := point(rover, base, baseOffset)
	if err != nil {
		log.Printf("point: %v\n", err)
		return
	}

	offset := s.pi.offset.Get().(*Attitude)
	s.pi.pan.Set(WrapAngle(att.Yaw + offset.Yaw))
	s.pi.tilt.Set(WrapAngle(att.Pitch + offset.Pitch))
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
