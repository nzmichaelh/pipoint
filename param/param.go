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

package param

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// Param is a value or a struct that has age and validity.  Updating a
// Param also fires an event.
type Param struct {
	Name    string
	value   interface{}
	updated time.Time
	params  *Params
	final   bool
}

// Ok return true if the value has been recently updated.
func (p *Param) Ok() bool {
	if p.final {
		return true
	}
	elapsed := time.Now().Sub(p.updated)
	return elapsed.Seconds() < 3
}

// Get returns the current value which may be nil.
func (p *Param) Get() interface{} {
	return p.value
}

// GetFloat64 returns the current value as a float64.
func (p *Param) GetFloat64() float64 {
	return p.value.(float64)
}

// GetInt returns the current value as an int.
func (p *Param) GetInt() int {
	return int(p.GetFloat64())
}

func asNumber(value interface{}) (float64, error) {
	if value == nil {
		return 0, errors.New("Is nil")
	}

	switch value.(type) {
	case float64:
		return value.(float64), nil
	case int:
		return float64(value.(int)), nil
	default:
		return 0, errors.New("Not a number")
	}
}

func isNumber(value interface{}) bool {
	_, err := asNumber(value)
	return err == nil
}

// Set the value, update validity, and notify listeners.
func (p *Param) Set(value interface{}) error {
	if p.value == nil {
		// OK, nothing set yet.
	} else if isNumber(p.value) && isNumber(value) {
		// Number -> number is fine.
	} else if reflect.TypeOf(value) != reflect.TypeOf(p.value) {
		return fmt.Errorf("Type of %v changed from %v to %v",
			p.Name, p.value, value)
	}

	if isNumber(value) {
		value, _ = asNumber(value)
	}

	p.value = value
	p.updated = time.Now()
	p.final = false
	p.params.updated(p)
	return nil
}

// Update tries to update the value.
func (p *Param) Update(value interface{}) (bool, error) {
	// Currently only handles numbers.
	if isNumber(value) && isNumber(p.value) {
		right, _ := asNumber(value)
		left, _ := asNumber(p.value)

		if left == right {
			return false, nil
		}
	}
	return true, p.Set(value)
}

// SetFloat64 tries to update the value as a float64.
func (p *Param) SetFloat64(value float64) error {
	return p.Set(value)
}

// SetInt tries to update the value as an int.
func (p *Param) SetInt(value int) error {
	return p.Set(value)
}

// UpdateInt tries to update the value as an int.
func (p *Param) UpdateInt(value int) (bool, error) {
	if p.GetInt() == value {
		return false, nil
	}
	return true, p.Set(value)
}

// Inc tries to increment the integer value.
func (p *Param) Inc() error {
	return p.SetInt(p.GetInt() + 1)
}

// Dec tries to decrement the integer value.
func (p *Param) Dec() error {
	return p.SetInt(p.GetInt() - 1)
}

// Finalise marks the param as always valid.
func (p *Param) Finalise() {
	p.final = true
}

// ValueVisitor is a callback for leaves in the parameter tree.
type ValueVisitor func(p *Param, path []string, value interface{})

func (p *Param) walk(visitor ValueVisitor, path []string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		p.walk(visitor, path, v.Elem())
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			p.walk(visitor, append(path, t.Field(i).Name), v.Field(i))
		}
	default:
		visitor(p, path, v.Interface())
	}
}

// Walk calls visitor on all leaf values of this parameter.
func (p *Param) Walk(visitor ValueVisitor) {
	p.walk(visitor, []string{p.Name}, reflect.ValueOf(p.Get()))
}
