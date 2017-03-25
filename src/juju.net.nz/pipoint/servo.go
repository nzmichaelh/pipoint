package pipoint

import (
	"math"
	"log"
)

// Parameters for a servo including limits.
type ServoParams struct {
	Pin int
	Span float64

	Min float64
	Max float64
}

// A servo on a pin with limits, demand, and actual position.
type Servo struct {
	params *Param
	sp     *Param
	pv     *Param
	pwm *PwmPin
}

func NewServo(name string, params *Params) *Servo {
	s := &Servo{
		params: params.NewWith(name, &ServoParams{
			Pin: -1,
			Min: 1.1,
			Max: 1.9,
			Span: math.Pi,
		}),
		sp: params.New(name + ".sp"),
		pv: params.New(name + ".pv"),
	}

	return s
}

func (s *Servo) Set(angle float64) {
	params := s.params.Get().(*ServoParams)

	s.sp.SetFloat64(angle)

	// Convert to pulse width.
	angle += math.Pi/2
	angle /= math.Pi

	ms := params.Min + angle*(params.Max - params.Min)
	if ms < 1 {
		ms = 1
	}
	if ms > 2 {
		ms = 2
	}
	s.pv.SetFloat64(ms)

	if params.Pin < 0 {
		return
	}
	if s.pwm == nil || params.Pin != s.pwm.Pin {
		if s.pwm != nil {
			s.pwm.SetEnable(0)
			s.pwm.UnExport()
		}
		s.pwm = &PwmPin{Chip: 0, Pin: params.Pin}
		
		var err error
		s.pwm.Export()
		err = s.pwm.SetEnable(0)
		if err == nil {
			err = s.pwm.SetPeriod(20e6)
		}
		if err == nil {
			s.pwm.SetDuty(int(ms*1e6))
		}
		if err == nil {
			err = s.pwm.SetEnable(1)
		}
		if err != nil {
			log.Printf("error: %v\n", err)
		}
	}
	if s.pwm != nil {
		s.pwm.SetDuty(int(ms*1e6))
	}
}
