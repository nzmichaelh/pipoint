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

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gobot.io/x/gobot/platforms/mqtt"
)

const (
	publishLimit = 0.2
)

// ParamMQTTBridge exposes parameters over MQTT.
type ParamMQTTBridge struct {
	params    *Params
	adaptor   *mqtt.Adaptor
	prefix    string
	limiter   *Limiter
	listening bool
}

// NewParamMQTTBridge creates a new two way connection between MQTT
// and the given params.
func NewParamMQTTBridge(params *Params, adaptor *mqtt.Adaptor, device string) *ParamMQTTBridge {
	if device == "" {
		device, _ = os.Hostname()
	}
	prefix := strings.Join([]string{device, params.Name}, "/")

	b := &ParamMQTTBridge{
		params:  params,
		adaptor: adaptor,
		prefix:  prefix,
		limiter: NewLimiter(),
	}
	changed := make(chan *Param, 10)
	params.Listen(changed)

	go func() {
		for {
			b.publish(<-changed)
		}
	}()
	return b
}

func (b *ParamMQTTBridge) publish(param *Param) {
	base := b.prefix + "/"

	param.Walk(func(p *Param, path []string, value interface{}) {
		name := base + strings.Join(path, "/")
		name = strings.ToLower(strings.Replace(name, ".", "/", -1))

		if b.limiter.Ok(name, publishLimit) {
			formatted := fmt.Sprintf("%v", value)
			b.adaptor.Publish(name, []byte(formatted))
		}
	})

	// TODO(michaelh): listen on connect.
	if !b.listening {
		b.listening = b.adaptor.On(b.prefix+"/#", b.recv)
	}
}

func (b *ParamMQTTBridge) recv(msg mqtt.Message) {
	topic := msg.Topic()
	if !strings.HasSuffix(topic, "/set") {
		return
	}
	if !strings.HasPrefix(topic, b.prefix) {
		return
	}

	// Convert into relative dotted form.
	parts := strings.Split(topic, "/")
	// Drop the device/ and /set.
	name := strings.Join(parts[1:len(parts)-1], ".")

	// Parse the data to a number or string.
	var next interface{}
	value := string(msg.Payload())
	fp, err := strconv.ParseFloat(value, 64)

	if err == nil {
		next = fp
	} else {
		next = value
	}

	b.params.WalkLeaves(func(p *Param, pname string, v reflect.Value) {
		if pname == name {
			if v.CanSet() {
				v.Set(reflect.ValueOf(next))
			} else {
				p.Set(next)
			}
		}
	})
}
