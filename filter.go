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

// Lowpass is a first order low pass filter.
type Lowpass struct {
	// The current filter output.
	Acc float64
}

// StepEx feeds v into the filter and returns the filtered value.  Tau
// is the filter coefficient, where 1.0 is no filtering and 0.0 is no
// pass through.
func (l *Lowpass) StepEx(v, tau float64) float64 {
	l.Acc = tau*v + (1-tau)*l.Acc
	return l.Acc
}

// LinPred is a velocity based linear predictive filter.
type LinPred struct {
	x       float64
	stamp   float64
	updated float64
	v       float64
}

// SetEx is called when a new measurement arrives.  x is the
// measurement, now is the local time this function was called, and
// stamp is the sensor specific time the measurement was made.
//
// The velocity is calculated based on sensor time.  The prediction is
// based on local time.
func (l *LinPred) SetEx(x, now, stamp float64) {
	if l.stamp == 0 {
		// First run
		l.v = 0
	} else {
		dt := stamp - l.stamp
		dx := x - l.x
		if dt <= 0 {
			l.v = 0
		} else {
			l.v = dx / dt
		}
	}
	l.x = x
	l.stamp = stamp
	l.updated = now
}

// GetEx returns the predicted measurement based on the given local
// time.
func (l *LinPred) GetEx(now float64) float64 {
	dt := now - l.updated
	if dt < 0 {
		dt = 0
	} else if dt > 2 {
		// Clamp if the value hasn't been updated recently.
		dt = 2
	}
	return l.x + l.v*dt
}
