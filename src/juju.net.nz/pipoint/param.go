package main

import (
	"log"
	"time"
)

type Param struct {
	name    string
	value   interface{}
	updated time.Time
	params  *Params
	final   bool
}

func (p *Param) Ok() bool {
	if p.final {
		return true
	}
	elapsed := time.Now().Sub(p.updated)
	return elapsed.Seconds() < 3
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

func (p *Param) Get() interface{} {
	return p.value
}

func (p *Param) Set(value interface{}) {
	log.Printf("set %v = %v\n", p.name, value)
	p.value = value
	p.updated = time.Now()
	p.final = false
	p.params.Updated(p)
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

func (p *Param) Final() {
	p.final = true
}

type ParamListener func(p *Param)

type Params struct {
	params    map[string]Param
	listeners []ParamListener
}

func (ps *Params) New(name string) *Param {
	p := &Param{
		name:   name,
		params: ps,
	}
	return p
}

func (ps *Params) NewWith(name string, value interface{}) *Param {
	p := &Param{
		name:   name,
		params: ps,
	}
	p.Set(value)
	return p
}

func (ps *Params) Updated(p *Param) {
	for _, l := range ps.listeners {
		l(p)
	}
}

func (ps *Params) Listen(l ParamListener) {
	ps.listeners = append(ps.listeners, l)
}
