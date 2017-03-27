package pipoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinPred(t *testing.T) {
	l := &LinPred{}

	// Returns initial value.
	l.FeedEx(5, 1)
	assert.InDelta(t, l.GetEx(1), 5, 0.001)
	// Initial value continues.
	assert.InDelta(t, l.GetEx(2), 5, 0.001)

	l.FeedEx(5.5, 2)
	// Initially returns the same value.
	assert.InDelta(t, l.GetEx(2), 5.5, 0.001)
	// Velocity is 0.5/s
	assert.InDelta(t, l.GetEx(2.5), 5.75, 0.001)
	assert.InDelta(t, l.GetEx(3.0), 6.0, 0.001)

	// Negative velocity also works.
	l.FeedEx(4, 3)
	assert.InDelta(t, l.GetEx(3), 4, 0.001)
	assert.InDelta(t, l.GetEx(4), 2.5, 0.001)
}
