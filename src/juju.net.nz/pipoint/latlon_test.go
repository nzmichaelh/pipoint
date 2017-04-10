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
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatLength(t *testing.T) {
	// Test vectors from https://en.wikipedia.org/wiki/Latitude
	assert.InDelta(t, LatLength(0*math.Pi/180), 110574, 1)
	assert.InDelta(t, LatLength(15*math.Pi/180), 110649, 1)
	assert.InDelta(t, LatLength(30*math.Pi/180), 110852, 1)
	assert.InDelta(t, LatLength(45*math.Pi/180), 111132, 1)
	assert.InDelta(t, LatLength(60*math.Pi/180), 111412, 1)
	assert.InDelta(t, LatLength(75*math.Pi/180), 111618, 1)
}

func TestLonLength(t *testing.T) {
	// Test vectors from https://en.wikipedia.org/wiki/Latitude
	assert.InDelta(t, LonLength(0*math.Pi/180), 111320, 1)
	assert.InDelta(t, LonLength(15*math.Pi/180), 107550, 1)
	assert.InDelta(t, LonLength(30*math.Pi/180), 96486, 1)
	assert.InDelta(t, LonLength(45*math.Pi/180), 78847, 1)
	assert.InDelta(t, LonLength(60*math.Pi/180), 55800, 1)
	assert.InDelta(t, LonLength(75*math.Pi/180), 28902, 1)
}

func ExampleLatLength() {
	fmt.Println(int(LatLength(46.8 * math.Pi / 180)))
	// Output: 111166
}
