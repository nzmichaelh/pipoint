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
