package pipoint

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPointEast(t *testing.T) {
	rover := &NEUPosition{North: 2000, East: 3100}
	base := &NEUPosition{North: 2000, East: 3000}

	at, err := point(rover, base, &NEUPosition{})
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, 0, 0.001)
	assert.InDelta(t, at.Yaw, math.Pi/2, 0.001)
}

func TestPointWest(t *testing.T) {
	rover := &NEUPosition{North: 2000, East: 2900}
	base := &NEUPosition{North: 2000, East: 3000}

	at, err := point(rover, base, &NEUPosition{})
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, 0, 0.001)
	assert.InDelta(t, at.Yaw, -math.Pi/2, 0.001)
}

func TestPointNorth(t *testing.T) {
	rover := &NEUPosition{North: 2100, East: 3000}
	base := &NEUPosition{North: 2000, East: 3000}

	at, err := point(rover, base, &NEUPosition{})
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, 0, 0.001)
	assert.InDelta(t, at.Yaw, 0, 0.001)
}

func TestPointSouth(t *testing.T) {
	rover := &NEUPosition{North: 1900, East: 3000}
	base := &NEUPosition{North: 2000, East: 3000}

	at, err := point(rover, base, &NEUPosition{})
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, 0, 0.001)
	assert.InDelta(t, at.Yaw, math.Pi, 0.001)
}

func TestPointUp(t *testing.T) {
	// ~45 degrees up.
	rover := &NEUPosition{North: 2000, East: 4000, Up: 1000}
	base := &NEUPosition{North: 2000, East: 3000, Up: 0}

	at, err := point(rover, base, &NEUPosition{})
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, math.Pi/4, 0.01)
	assert.InDelta(t, at.Yaw, math.Pi/2, 0.001)
}

func TestPointDown(t *testing.T) {
	// ~45 degrees down.
	rover := &NEUPosition{North: 2000, East: 4000, Up: 0}
	base := &NEUPosition{North: 2000, East: 3000, Up: 1000}

	at, err := point(rover, base, &NEUPosition{})
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, -math.Pi/4, 0.01)
	assert.InDelta(t, at.Yaw, math.Pi/2, 0.001)
}
