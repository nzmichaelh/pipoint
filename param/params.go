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
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/spf13/viper"
)

// ParamChannel passes changes to a Param.
type ParamChannel chan *Param

// Params is a group of parameters that can be listened to.
type Params struct {
	Name      string
	viper     *viper.Viper
	params    []*Param
	listeners []ParamChannel
}

// NewParams creates a new group of parameters with the given root name.
func NewParams(name string) *Params {
	ps := &Params{
		Name:  name,
		viper: viper.New(),
	}

	ps.viper.SetConfigName(name)
	ps.viper.AddConfigPath(".")
	ps.viper.ReadInConfig()
	ps.viper.WatchConfig()
	ps.viper.OnConfigChange(func(e fsnotify.Event) { ps.Load() })

	return ps
}

// New creates a new Param in this group.  The Param is uninitialised
// and invalid.
func (ps *Params) New(name string) *Param {
	p := &Param{
		Name:   name,
		params: ps,
	}
	ps.params = append(ps.params, p)
	return p
}

// NewNum create a new number param in this group.  The Param is zero
// and invalid.
func (ps *Params) NewNum(name string) *Param {
	p := &Param{
		Name:   name,
		value:  0.0,
		params: ps,
	}
	ps.params = append(ps.params, p)
	return p
}

// NewWith create a new, valid Param in this group using the given
// value.
func (ps *Params) NewWith(name string, value interface{}) *Param {
	p := &Param{
		Name:   name,
		params: ps,
	}
	p.Set(value)
	ps.params = append(ps.params, p)
	return p
}

func (ps *Params) updated(p *Param) {
	for _, l := range ps.listeners {
		l <- p
	}
}

// Listen to changes on any parameter in this group.
func (ps *Params) Listen(l ParamChannel) {
	ps.listeners = append(ps.listeners, l)
}

// LeafVisitor is called on every param value.
type LeafVisitor func(p *Param, name string, value reflect.Value)

// WalkLeaves calls the given visitor on every leaf value.
func (ps *Params) WalkLeaves(visitor LeafVisitor) {
	for _, p := range ps.params {
		ps.visitOne(p, visitor, []string{ps.Name, p.Name}, reflect.ValueOf(p.Get()))
	}
}

func makeName(path []string) string {
	name := strings.Join(path, ".")
	return strings.ToLower(name)
}

func (ps *Params) visitOne(p *Param, visitor LeafVisitor, path []string, v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		ps.visitOne(p, visitor, path, v.Elem())
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			ps.visitOne(p, visitor, append(path, t.Field(i).Name), v.Field(i))
		}
	default:
		visitor(p, makeName(path), v)
	}
}

func format(name string, v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return ""
	case reflect.String:
		if v.String() != "" {
			return fmt.Sprintf("%s{value=\"%s\"} 1", name, v)
		}
		return ""
	case reflect.Float64, reflect.Float32:
		vf := v.Float()
		if float64(int(vf)) == vf {
			return fmt.Sprintf("%s %v", name, int(vf))
		}
		return fmt.Sprintf("%s %v", name, v)
	default:
		return fmt.Sprintf("%s %v", name, v)
	}
}

// Metrics generates Prometheus/Borgmon compatible metrics from all
// params.
func (ps *Params) Metrics(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer

	ps.WalkLeaves(func(p *Param, name string, v reflect.Value) {
		name = strings.Replace(name, ".", "_", -1)
		f := format(name, v)
		if f != "" {
			buf.WriteString(f + "\n")
		}
	})

	enc := expfmt.NewEncoder(&buf, expfmt.FmtText)
	mfs, err := prometheus.DefaultGatherer.Gather()

	if err == nil {
		for _, mf := range mfs {
			enc.Encode(mf)
		}
	}

	header := w.Header()
	header.Set("Content-Type", string(expfmt.FmtText))
	w.Write(buf.Bytes())
}

// Load fetches the supplied values from Viper and updates all
// matching params.
func (ps *Params) Load() {
	ps.WalkLeaves(func(p *Param, name string, v reflect.Value) {
		if !ps.viper.IsSet(name) {
			return
		}
		next := ps.viper.Get(name)
		if v.CanSet() {
			v.Set(reflect.ValueOf(next))
		} else {
			p.Update(next)
		}
	})
}
