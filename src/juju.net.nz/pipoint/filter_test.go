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
