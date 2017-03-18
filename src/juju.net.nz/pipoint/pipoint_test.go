package pipoint

import (
	"testing"
	"math"
)

func TestPointEast(t *testing.T) {
	// ~7 km to the east.
	rover := &Position{46.8, 8.3, 0}
	base := &Position{46.8, 8.2, 0}

	at, err := NewPiPoint().point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	checkClose(t, at.Pitch, 0, 0.001)
	checkClose(t, at.Yaw, math.Pi/2, 0.001)
}

func TestPointWest(t *testing.T) {
	// ~7 km to the west.
	rover := &Position{46.8, 8.1, 0}
	base := &Position{46.8, 8.2, 0}

	at, err := NewPiPoint().point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	checkClose(t, at.Pitch, 0, 0.001)
	checkClose(t, at.Yaw, -math.Pi/2, 0.001)
}

func TestPointNorth(t *testing.T) {
	// ~10 km to the north.
	rover := &Position{46.9, 8.2, 0}
	base := &Position{46.8, 8.2, 0}

	at, err := NewPiPoint().point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	checkClose(t, at.Pitch, 0, 0.001)
	checkClose(t, at.Yaw, 0, 0.001)
}

func TestPointSouth(t *testing.T) {
	// ~10 km to the south.
	rover := &Position{46.7, 8.2, 0}
	base := &Position{46.8, 8.2, 0}

	at, err := NewPiPoint().point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	checkClose(t, at.Pitch, 0, 0.001)
	checkClose(t, at.Yaw, math.Pi, 0.001)
}

func TestPointUp(t *testing.T) {
	// ~50 degrees up.
	rover := &Position{46.8, 8.21, 1000}
	base := &Position{46.8, 8.2, 0}

	at, err := NewPiPoint().point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	checkClose(t, at.Pitch, AsRad(53), 0.01)
	checkClose(t, at.Yaw, math.Pi/2, 0.001)
}

func TestPointDown(t *testing.T) {
	// ~50 degrees down.
	rover := &Position{46.8, 8.21, 0}
	base := &Position{46.8, 8.2, 1000}

	at, err := NewPiPoint().point(rover, base)
	if err != nil {
		t.Errorf("%#v", err)
	}
	checkClose(t, at.Pitch, AsRad(-53), 0.01)
	checkClose(t, at.Yaw, math.Pi/2, 0.001)
}
