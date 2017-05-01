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

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"math"
	"regexp"
	"strings"
	"time"
)

// Scale from one range to another.
func Scale(v, min0, max0, min1, max1 float64) float64 {
	v -= min0
	v /= (max0 - min0)
	v *= (max1 - min1)
	v += min1

	return v
}

// WrapAngle wraps an angle in radians to [-pi, pi)
func WrapAngle(v float64) float64 {
	// TODO: add tests.
	for v >= math.Pi {
		v -= math.Pi * 2
	}
	for v < -math.Pi {
		v += math.Pi * 2
	}
	return v
}

var nowOverride *float64

// Now returns the current system time in seconds.
func Now() float64 {
	if nowOverride != nil {
		return *nowOverride
	}
	// Pulled out so it can be mocked.
	return float64(time.Now().UnixNano()) * 1e-9
}

// OverrideNow allows overriding Now() for testing.
func OverrideNow(now float64) {
	nowOverride = &now
}

// NormText returns a short, unique version of the string.
func NormText(text string) string {
	re := regexp.MustCompile("\\W")

	hashed := sha1.New()
	io.WriteString(hashed, text)
	summary := hex.EncodeToString(hashed.Sum(nil))[:4]

	norm := text
	if len(norm) > 20 {
		norm = norm[:20]
	}
	norm = re.ReplaceAllString(text, "_")
	return strings.ToLower(norm + "-" + summary)
}

func ServoToScale(ms uint16) float64 {
	v := float64(ms)
	return (v - 1500) / 500
}

func ScaleToPos(v float64) int {
	switch {
	case v < -0.5:
		return -2
	case v < -0.2:
		return -1
	case v < 0.2:
		return 0
	case v < 0.5:
		return 1
	default:
		return 2
	}
}
