package pipoint

import (
	"fmt"
	"reflect"
	"time"
)

// Param is a value or a struct that has age and validity.  Updating a
// Param also fires an event.
type Param struct {
	name    string
	value   interface{}
	updated time.Time
	params  *Params
	final   bool
}

// Return true if the value has been recently updated.
func (p *Param) Ok() bool {
	if p.final {
		return true
	}
	elapsed := time.Now().Sub(p.updated)
	return elapsed.Seconds() < 3
}

func (p *Param) Get() interface{} {
	return p.value
}

func (p *Param) GetFloat64() float64 {
	if p.value == nil {
		return 0
	}
	return p.value.(float64)
}

func (p *Param) GetInt() int {
	if p.value == nil {
		return 0
	}
	return p.value.(int)
}

// Set the value, update validity, and notify listeners.
func (p *Param) Set(value interface{}) {
	if p.value != nil && reflect.TypeOf(value) !=
		reflect.TypeOf(p.value) {
		panic(fmt.Sprintf("Type of %v changed from %v to %v",
			p.name, p.value, value))
	}
	p.value = value
	p.updated = time.Now()
	p.final = false
	p.params.updated(p)
}

func (p *Param) SetFloat64(value float64) {
	p.Set(value)
}

func (p *Param) SetInt(value int) {
	p.Set(value)
}

func (p *Param) Inc() {
	p.SetInt(p.GetInt() + 1)
}

// Mark the param as always valid.
func (p *Param) Final() {
	p.final = true
}
