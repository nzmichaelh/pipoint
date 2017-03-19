package pipoint

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

	ps.Listen(func(p *Param) {
		hits += 1
		val = p
	})

	p.SetFloat64(17)
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
	assert.Panics(t, func() {
		p.Set(1.0)
	})
}
