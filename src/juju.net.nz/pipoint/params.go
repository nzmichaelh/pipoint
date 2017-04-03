package pipoint

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

type ParamListener func(p *Param)

// Params is a group of parameters that can be listened to.
type Params struct {
	name      string
	viper     *viper.Viper
	params    []*Param
	listeners []ParamListener
}

func NewParams(name string) *Params {
	ps := &Params{
		name:  name,
		viper: viper.New(),
	}

	ps.viper.SetConfigName(name)
	ps.viper.AddConfigPath(".")
	ps.viper.ReadInConfig()
	ps.viper.WatchConfig()
	ps.viper.OnConfigChange(func(e fsnotify.Event) { ps.Load() })

	return ps
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

type LeafVisitor func(p *Param, path []string, value reflect.Value)

func (ps *Params) WalkLeaves(visitor LeafVisitor) {
	for _, p := range ps.params {
		ps.visitOne(p, visitor, []string{ps.name, p.name}, reflect.ValueOf(p.Get()))
	}
}

func makeName(path []string, sep string) string {
	name := strings.Join(path, ".")
	return strings.Replace(strings.ToLower(name), ".", sep, -1)
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
		visitor(p, path, v)
	}
}

// Respond with the Prometheus and Params as metrics.
func (ps *Params) Metrics(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer

	ps.WalkLeaves(func(p *Param, path []string, v reflect.Value) {
		name := makeName(path, "_")
		switch v.Kind() {
		case reflect.Invalid:
			break
		case reflect.String:
			if v.String() != "" {
				buf.WriteString(fmt.Sprintf("%s{value=\"%s\"} 1\n", name, v))
			}
		default:
			buf.WriteString(fmt.Sprintf("%s %v\n", name, v))
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

func (ps *Params) Load() {
	ps.WalkLeaves(func(p *Param, path []string, v reflect.Value) {
		name := makeName(path, ".")
		if !ps.viper.IsSet(name) {
			return
		}
		next := ps.viper.Get(name)
		if v.CanSet() {
			v.Set(reflect.ValueOf(next))
		} else {
			p.Set(next)
		}
	})
}
