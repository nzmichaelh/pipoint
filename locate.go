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

// LocateState runs initially to locate the base unit.
type LocateState struct {
	name string
	pi   *PiPoint
}

func (s *LocateState) Name() string {
	return "Locate"
}

// Update is called when a param is updated.
func (s *LocateState) Update(param *param.Param) {
	switch param {
	case s.pi.neu:
		s.pi.rover.Set(param.Get())
		s.pi.base.Set(param.Get())
		s.pi.base.Finalise()
	case s.pi.attitude:
		s.pi.offset.Set(&Attitude{
			Yaw: s.pi.attitude.Get().(*Attitude).Yaw,
		})
	case s.pi.mark:
		s.pi.state.Inc()
	}
}
