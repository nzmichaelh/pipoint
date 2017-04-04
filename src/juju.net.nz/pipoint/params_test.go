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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParamsFormat(t *testing.T) {
	ps := Params{}

	// Whole numbers print as ints.
	p := ps.NewWith("foo", 17.0)
	assert.Equal(t, format("foo", reflect.ValueOf(p.Get())), "foo 17")

	p.SetFloat64(17.1)
	assert.Equal(t, format("foo", reflect.ValueOf(p.Get())), "foo 17.1")

	// Strings show only if set.
	p2 := ps.NewWith("bar", "texty")
	assert.Equal(t, format("foo", reflect.ValueOf(p2.Get())), "foo{value=\"texty\"} 1")

	p2.Set("")
	assert.Equal(t, format("foo", reflect.ValueOf(p2.Get())), "")
}
