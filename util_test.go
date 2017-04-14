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

func TestNormText(t *testing.T) {
	assert.Equal(t, NormText("Foo"),
		"foo-201a")
	assert.Equal(t, NormText("Base location set"),
		"base_location_set-4fe5")
	assert.Equal(t, NormText("Lots%of--crazy"),
		"lots_of__crazy-4921")

	
	// No collisions.
	assert.Equal(t, NormText("Rover online"),
		"rover_online-f40c")
	assert.Equal(t, NormText("Rover offline"),
		"rover_offline-5951")
}
