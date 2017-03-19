package pipoint

import (
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

type ParamListener func(p *Param)

// Params is a group of parameters that can be listened to.
type Params struct {
	params    map[string]Param
	listeners []ParamListener
}

// Create a new Param in this group.  The Param is uninitialised and
// invalid.
func (ps *Params) New(name string) *Param {
	p := &Param{
		name:   name,
		params: ps,
	}
	return p
}

// Create a new, valid Param in this group using the given value.
func (ps *Params) NewWith(name string, value interface{}) *Param {
	p := &Param{
		name:   name,
		params: ps,
	}
	p.Set(value)
	return p
}

func (ps *Params) updated(p *Param) {
	for _, l := range ps.listeners {
		l(p)
	}
}

// Listen to changes on any parameter in this group.
func (ps *Params) Listen(l ParamListener) {
	ps.listeners = append(ps.listeners, l)
}
