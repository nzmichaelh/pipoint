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
	"testing"

	"juju.net.nz/x/pipoint/param"
	"juju.net.nz/x/pipoint/util"

	"github.com/stretchr/testify/assert"
)

func TestServo(t *testing.T) {
	p := &param.Params{}
	s := NewServo("pan", p)

	util.OverrideNow(1)

	assert.InDelta(t, s.sp.GetFloat64(), 0, 0.001, "Starts at zero")

	s.Set(-math.Pi * 0.5)
	s.Tick()
	assert.InDelta(t, s.pv.GetFloat64(), 1.1, 0.01)
	s.Set(+math.Pi * 0.5)
	s.Tick()
	assert.InDelta(t, s.pv.GetFloat64(), 1.9, 0.01)

	// Stays within limits
	s.Set(-math.Pi * 0.7)
	s.Tick()
	assert.InDelta(t, s.pv.GetFloat64(), 1.0, 0.01)
	s.Set(+math.Pi * 0.7)
	s.Tick()
	assert.InDelta(t, s.pv.GetFloat64(), 2.0, 0.01)
}
