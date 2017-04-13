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

func TestParam(t *testing.T) {
	ps := Params{}

	p := ps.New("foo")

	// Never set, so not OK.
	assert.False(t, p.Ok())

	// Can set and get.
	p.SetFloat64(3)
	assert.Equal(t, p.GetFloat64(), 3.0)
}

func TestParamListen(t *testing.T) {
	ps := Params{}
	p := ps.New("foo")

	hits := 0
	var val *Param

	ch := make(ParamChannel, 10)
	ps.Listen(ch)

	p.SetFloat64(17)

	close(ch)

	for p := range(ch) {
		hits++
		val = p
	}

	assert.Equal(t, val, p)
	assert.Equal(t, hits, 1)
}

type TestParamStructT struct {
	A int
	B int
	C float64
}

func TestParamStruct(t *testing.T) {
	ps := Params{}
	p := ps.NewWith("blob", &TestParamStructT{1, 2, 3})

	v := p.Get().(*TestParamStructT)

	assert.Equal(t, v.A, 1)
	assert.Equal(t, v.B, 2)
	assert.Equal(t, v.C, 3.0)

	// Trying to set a different type causes an error.
	assert.Error(t, p.Set(1.0))
}

func TestParamNumber(t *testing.T) {
	ps := Params{}
	p := ps.NewNum("blob")

	// Defaults to zero.
	assert.Equal(t, p.GetFloat64(), 0.0)

	// Can set to ints or floats
	p.SetFloat64(3.5)
	assert.Equal(t, p.GetFloat64(), 3.5)
	p.SetInt(17)
	assert.Equal(t, p.GetInt(), 17)

	// Setting to a non-number casues an error.
	assert.Error(t, p.Set(&TestParamStructT{}))
}
