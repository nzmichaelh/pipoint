package pipoint

import (
	"math"
)

type CycleState struct {
	pi    *PiPoint
	cycle float64
}

func (s *CycleState) Update(param *Param) {
	switch param {
	case s.pi.tick:
		s.cycle += 0.02
		s.pi.pan.Set(Scale(math.Cos(s.cycle), -1, 1, -math.Pi/2, math.Pi/2))
		s.pi.tilt.Set(Scale(math.Sin(s.cycle), -1, 1, -math.Pi/2, 0))
	}
}
