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

package util

// Limiter is a simple keyed rate limiter.
type Limiter struct {
	stamp map[string]float64
}

// NewLimiter returns a new, initialised limiter.
func NewLimiter() *Limiter {
	return &Limiter{
		stamp: make(map[string]float64),
	}
}

// Ok returns true if at least dt seconds have passed since the last
// time this function returned true.
func (l *Limiter) Ok(key string, dt float64) bool {
	now := Now()

	stamp, ok := l.stamp[key]

	if !ok {
		l.stamp[key] = now
		return true
	}

	elapsed := now - stamp

	if elapsed < dt {
		return false
	}

	if elapsed < 2*dt {
		l.stamp[key] = stamp + dt
	} else {
		l.stamp[key] = now
	}

	return true
}
