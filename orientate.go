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
	"juju.nz/x/pipoint/param"
)

// OrientateState runs to set the pan orientation.
type OrientateState struct {
	name string
	pi   *PiPoint
}

func (s *OrientateState) Name() string {
	return "Orientate"
}

// Update is called when a param is updated.
func (s *OrientateState) Update(param *param.Param) {
	switch param {
	case s.pi.neu:
		s.pi.rover.Set(param.Get())
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
	offset := s.pi.baseOffset.Get().(*NEUPosition)

	att, err := point(rover, base, offset)

	if err != nil {
		return
	}

	current := s.pi.offset.Get().(*Attitude)
	s.pi.offset.Set(&Attitude{
		Yaw:   -att.Yaw,
		Pitch: current.Pitch,
	})
}
