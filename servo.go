// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package pipoint

import (
	"log"
	"math"
)

// ServoParams holds the parameters for a servo including limits.
type ServoParams struct {
	Pin  int
	Span float64

	Min  float64
	Max  float64
	Low  float64
	High float64

	Tau float64
}

// Servo is a servo on a pin with limits, demand, and actual
// position.
type Servo struct {
	params *Param
	sp     *Param
	pv     *Param
	pwm    *PwmPin
	filter *Lowpass
}

// NewServo creates a new servo with params on the given tree.
func NewServo(name string, params *Params) *Servo {
	s := &Servo{
		params: params.NewWith(name, &ServoParams{
			Pin:  -1,
			Min:  1.0,
			Max:  2.0,
			Low:  1.1,
			High: 1.9,
			Span: math.Pi,
			Tau:  1.0,
		}),
		sp:     params.NewNum(name + ".sp"),
		pv:     params.NewNum(name + ".pv"),
		filter: &Lowpass{},
	}

	return s
}

// Set updates the target angle in radians.  The servo is actually
// updated on calling Tick().
func (s *Servo) Set(angle float64) {
	s.sp.SetFloat64(angle)
}

// Tick updates the servo output based on demand.  Call every ~20 ms.
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
		// PWM pin has been set or changed.  Update.
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
