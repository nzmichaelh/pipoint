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
