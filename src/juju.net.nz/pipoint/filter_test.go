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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLowpass(t *testing.T) {
	l := &Lowpass{Tau: 0.3}

	// Moves towards the input in diminishing steps.
	assert.InDelta(t, l.Step(3), 0.9, 0.1)
	assert.InDelta(t, l.Step(3), 1.6, 0.1)
	assert.InDelta(t, l.Step(3), 2.0, 0.1)
	assert.InDelta(t, l.Step(3), 2.3, 0.1)
	assert.InDelta(t, l.Step(3), 2.5, 0.1)
}
