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

func TestLinPred(t *testing.T) {
	l := &LinPred{}

	// Returns initial value.
	l.SetEx(5, 1, 1)
	assert.InDelta(t, l.GetEx(1), 5, 0.001)
	// Initial value continues.
	assert.InDelta(t, l.GetEx(2), 5, 0.001)

	l.SetEx(5.5, 2, 2)
	// Initially returns the same value.
	assert.InDelta(t, l.GetEx(2), 5.5, 0.001)
	// Velocity is 0.5/s
	assert.InDelta(t, l.GetEx(2.5), 5.75, 0.001)
	assert.InDelta(t, l.GetEx(3.0), 6.0, 0.001)

	// Negative velocity also works.
	l.SetEx(4, 3, 3)
	assert.InDelta(t, l.GetEx(3), 4, 0.001)
	assert.InDelta(t, l.GetEx(4), 2.5, 0.001)
}
