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

type Lowpass struct {
	Tau float64
	Acc float64
}

func (l *Lowpass) Step(v float64) float64 {
	return l.StepEx(v, l.Tau)
}

func (l *Lowpass) StepEx(v, tau float64) float64 {
	l.Acc = tau*v + (1-tau)*l.Acc
	return l.Acc
}
