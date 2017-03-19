package pipoint

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestServo(t *testing.T) {
	p := &Params{}
	s := NewServo("pan", p)

	assert.InDelta(t, s.sp.GetFloat64(), 0, 0.001, "Starts at zero")

	s.Set(-math.Pi / 7)
	assert.InDelta(t, s.pv.GetFloat64(), -math.Pi/7, 0.01)
	s.Set(math.Pi / 3)
	assert.InDelta(t, s.pv.GetFloat64(), math.Pi/3, 0.01)

	// Stays within limits
	s.Set(-math.Pi * 2 / 3)
	assert.InDelta(t, s.pv.GetFloat64(), -math.Pi/2, 0.01)
	s.Set(+math.Pi * 2 / 3)
	assert.InDelta(t, s.pv.GetFloat64(), math.Pi/2, 0.01)
}
