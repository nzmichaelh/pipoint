package pipoint

import (
	"math"
)

// Parameters for a servo including limits.
type ServoParams struct {
	Pin string
	Min float64
	Max float64
	Mid float64
}

// A servo on a pin with limits, demand, and actual position.
type Servo struct {
	params *Param
	sp     *Param
	pv     *Param
}

func NewServo(name string, params *Params) *Servo {
	s := &Servo{
		params: params.NewWith(name, &ServoParams{
			"",
			-math.Pi / 2,
			math.Pi / 2,
			0,
		}),
		sp: params.New(name + ".sp"),
		pv: params.New(name + ".pv"),
	}

	return s
}

func (s *Servo) Set(angle float64) {
	params := s.params.Get().(*ServoParams)

	s.sp.SetFloat64(angle)

	angle += params.Mid
	angle = math.Min(angle, params.Max)
	angle = math.Max(angle, params.Min)

	// Move from -pi/2..pi/2 to the gobot format.
	demand := int(AsDeg(angle) + 90 + 0.5)

	if demand < 0 {
		demand = 0
	}
	if demand > 250 {
		demand = 250
	}
	s.pv.SetFloat64(AsRad(float64(demand - 90)))
}
