package pipoint

import (
	"bytes"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"net/http"
	"reflect"
	"strings"
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

type ParamListener func(p *Param)

// Params is a group of parameters that can be listened to.
type Params struct {
	Name      string
	params    []*Param
	listeners []ParamListener
}

// Create a new Param in this group.  The Param is uninitialised and
// invalid.
func (ps *Params) New(name string) *Param {
	p := &Param{
		name:   name,
		params: ps,
	}
	ps.params = append(ps.params, p)
	return p
}

// Create a new, valid Param in this group using the given value.
func (ps *Params) NewWith(name string, value interface{}) *Param {
	p := &Param{
		name:   name,
		params: ps,
	}
	p.Set(value)
	ps.params = append(ps.params, p)
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

func (ps *Params) Walk(visitor func(*Param)) {
	for _, p := range ps.params {
		visitor(p)
	}
}

func makeName(path []string) string {
	name := strings.Join(path, ".")
	return strings.Replace(strings.ToLower(name), ".", "_", -1)
}

func (ps *Params) visitOne(w *bytes.Buffer, path []string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		ps.visitOne(w, path, v.Elem())
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			ps.visitOne(w, append(path, t.Field(i).Name), v.Field(i))
		}
	case reflect.String:
		w.WriteString(fmt.Sprintf("%s{value=\"%s\"} 1\n", makeName(path), v))
	default:
		w.WriteString(fmt.Sprintf("%s %v\n", makeName(path), v))
	}
}

// Respond with the Prometheus and Params as metrics.
func (ps *Params) Metrics(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	enc := expfmt.NewEncoder(&buf, expfmt.FmtText)
	mfs, err := prometheus.DefaultGatherer.Gather()

	if err == nil {
		for _, mf := range mfs {
			enc.Encode(mf)
		}
	}

	header := w.Header()
	header.Set("Content-Type", string(expfmt.FmtText))

	ps.Walk(func(p *Param) {
		if p.Get() != nil {
			ps.visitOne(&buf, []string{"pipoint", p.name}, reflect.ValueOf(p.Get()))
		}
	})

	w.Write(buf.Bytes())
}
