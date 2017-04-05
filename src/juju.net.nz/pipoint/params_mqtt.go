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
	"strings"

	"gobot.io/x/gobot/platforms/mqtt"
)

const (
	publishLimit = 0.2
)

// ParamMQTTBridge exposes parameters over MQTT.
type ParamMQTTBridge struct {
	adaptor *mqtt.Adaptor
	prefix  string
	limiter *Limiter
}

func NewParamMQTTBridge(params *Params, adaptor *mqtt.Adaptor, device string) *ParamMQTTBridge {
	if device == "" {
		device, _ = os.Hostname()
	}
	prefix := strings.Join([]string{device, params.Name}, "/")

	b := &ParamMQTTBridge{
		adaptor: adaptor,
		prefix:  prefix,
		limiter: NewLimiter(),
	}
	params.Listen(b.updated)
	return b
}

func (b *ParamMQTTBridge) updated(param *Param) {
	base := b.prefix + "/"

	param.Walk(func(p *Param, path []string, value interface{}) {
		name := base + strings.Join(path, "/")
		name = strings.ToLower(strings.Replace(name, ".", "/", -1))

		if b.limiter.Ok(name, publishLimit) {
			formatted := fmt.Sprintf("%v", value)
			b.adaptor.Publish(name, []byte(formatted))
		}
	})
}
