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

	"juju.nz/x/pipoint/param"
	"juju.nz/x/pipoint/util"
)

// CycleState cycles the pan/tilt in a circle.
type CycleState struct {
	pi    *PiPoint
	cycle float64
}

func (s *CycleState) Name() string {
	return "Cycle"
}

// Update reacts to changes in parameters.
func (s *CycleState) Update(param *param.Param) {
	switch param {
	case s.pi.tick:
		s.cycle += 0.02
		s.pi.pan.Set(util.Scale(math.Cos(s.cycle), -1, 1, -math.Pi/2, math.Pi/2))
		s.pi.tilt.Set(util.Scale(math.Sin(s.cycle), -1, 1, -math.Pi/2, 0))
	}
}
