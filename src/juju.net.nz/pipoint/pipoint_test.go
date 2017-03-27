package pipoint

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestPointEast(t *testing.T) {
	// ~7 km to the east.
	rover := &Position{Lat: 46.8, Lon: 8.3}
	base := &Position{Lat: 46.8, Lon: 8.2}

	at, err := point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, 0, 0.001)
	assert.InDelta(t, at.Yaw, math.Pi/2, 0.001)
}

func TestPointWest(t *testing.T) {
	// ~7 km to the west.
	rover := &Position{Lat: 46.8, Lon: 8.1}
	base := &Position{Lat: 46.8, Lon: 8.2}

	at, err := point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, 0, 0.001)
	assert.InDelta(t, at.Yaw, -math.Pi/2, 0.001)
}

func TestPointNorth(t *testing.T) {
	// ~10 km to the north.
	rover := &Position{Lat:46.9, Lon: 8.2}
	base := &Position{Lat:46.8, Lon: 8.2}

	at, err := point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, 0, 0.001)
	assert.InDelta(t, at.Yaw, 0, 0.001)
}

func TestPointSouth(t *testing.T) {
	// ~10 km to the south.
	rover := &Position{Lat:46.7, Lon: 8.2}
	base := &Position{Lat:46.8, Lon: 8.2}

	at, err := point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, 0, 0.001)
	assert.InDelta(t, at.Yaw, math.Pi, 0.001)
}

func TestPointUp(t *testing.T) {
	// ~50 degrees up.
	rover := &Position{Lat:46.8, Lon: 8.21, Alt: 1000}
	base := &Position{Lat:46.8, Lon: 8.2, Alt: 0}

	at, err := point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, AsRad(53), 0.01)
	assert.InDelta(t, at.Yaw, math.Pi/2, 0.001)
}

func TestPointDown(t *testing.T) {
	// ~50 degrees down.
	rover := &Position{Lat:46.8, Lon: 8.21, Alt: 0}
	base := &Position{Lat:46.8, Lon: 8.2, Alt: 1000}

	at, err := point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	assert.InDelta(t, at.Pitch, AsRad(-53), 0.01)
	assert.InDelta(t, at.Yaw, math.Pi/2, 0.001)
}
