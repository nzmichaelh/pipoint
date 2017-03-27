package pipoint

import (
	"log"
	"math"
)

// Parameters for a servo including limits.
type ServoParams struct {
	Pin  int
	Span float64

	Min  float64
	Max  float64
	Low  float64
	High float64

	Tau  float64
}

// A servo on a pin with limits, demand, and actual position.
type Servo struct {
	params *Param
	sp     *Param
	pv     *Param
	pwm    *PwmPin
	filter *Lowpass
}

func NewServo(name string, params *Params) *Servo {
	s := &Servo{
		params: params.NewWith(name, &ServoParams{
			Pin:  -1,
			Min:  1.0,
			Max:  2.0,
			Low:  1.1,
			High: 1.9,
			Span: math.Pi,
			Tau: 1.0,
		}),
		sp: params.New(name + ".sp"),
		pv: params.New(name + ".pv"),
		filter: &Lowpass{},
	}

	return s
}

func (s *Servo) Set(angle float64) {
	s.sp.SetFloat64(angle)
}

func (s *Servo) Tick() {
	params := s.params.Get().(*ServoParams)

	angle := s.sp.GetFloat64()
	angle = s.filter.StepEx(angle, params.Tau)

	// Convert to pulse width.
	angle += math.Pi / 2

	ms := Scale(angle, 0, params.Span, params.Low, params.High)
	ms = math.Min(params.Max, math.Max(params.Min, ms))
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
			s.pwm.SetDuty(int(ms * 1e6))
		}
		if err == nil {
			err = s.pwm.SetEnable(1)
		}
		if err != nil {
			log.Printf("error: %v\n", err)
		}
	}
	if s.pwm != nil {
		s.pwm.SetDuty(int(ms * 1e6))
	}
}
