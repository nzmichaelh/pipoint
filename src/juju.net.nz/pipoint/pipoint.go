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
	"log"
	"time"

	common "gobot.io/x/gobot/platforms/mavlink/common"
	"gobot.io/x/gobot/platforms/mqtt"
)

const (
	dt = time.Millisecond * 20
)

type State interface {
	Update(param *Param)
}

// An automatic, GPS based system that points a camera at the rover.
type PiPoint struct {
	Params *Params

	state *Param

	tick       *Param
	messages   *Param
	heartbeats *Param
	heartbeat  *Param
	attitude   *Param
	gps        *Param
	neu        *Param
	baseOffset *Param
	pred       *Param
	rover      *Param
	base       *Param
	sysStatus  *Param

	sp     *Param
	offset *Param

	pan  *Servo
	tilt *Servo

	latPred *LinPred
	lonPred *LinPred
	altPred *LinPred

	states []State

	elog *EventLogger
	log  *log.Logger

	param ParamChannel
}

// Create a new camera pointer.
func NewPiPoint() *PiPoint {
	p := &PiPoint{
		Params:  NewParams("pipoint"),
		latPred: &LinPred{},
		lonPred: &LinPred{},
		altPred: &LinPred{},
		elog:    NewEventLogger("pipoint"),
		param:   make(ParamChannel, 10),
	}

	p.log = p.elog.logger

	p.states = []State{
		&LocateState{pi: p},
		&RunState{pi: p},
		&CycleState{pi: p},
	}

	p.tick = p.Params.NewNum("tick")
	p.messages = p.Params.NewNum("rover.messages")

	p.state = p.Params.NewNum("state")
	p.heartbeat = p.Params.NewWith("heartbeat", &common.Heartbeat{})
	p.heartbeats = p.Params.NewNum("heartbeat")

	p.gps = p.Params.New("gps")
	p.neu = p.Params.New("position")
	p.pred = p.Params.New("pred")

	p.attitude = p.Params.New("rover.attitude")
	p.rover = p.Params.New("rover.position")
	p.base = p.Params.New("base.position")
	p.baseOffset = p.Params.NewWith("base.offset", &NEUPosition{})

	p.sysStatus = p.Params.New("rover.status")

	p.sp = p.Params.NewWith("pantilt.sp", &Attitude{})
	p.offset = p.Params.NewWith("pantilt.offset", &Attitude{})

	p.pan = NewServo("pantilt.pan", p.Params)
	p.tilt = NewServo("pantilt.tilt", p.Params)

	p.Params.Listen(p.param)
	p.Params.Load()
	return p
}

func (pi *PiPoint) AddMQTT(mqtt *mqtt.Adaptor) {
	NewParamMQTTBridge(pi.Params, mqtt, "")
}

func (pi *PiPoint) Run() {
	tick := time.NewTicker(dt)

	for {
		select {
		case param := <-pi.param:
			pi.update(param)
		case <-tick.C:
			pi.ticked()
		}
	}
}

func (pi *PiPoint) ticked() {
	now := Now()
	pi.tick.SetFloat64(now)

	pred := &Position{
		Time: now,
		Lat:  pi.latPred.GetEx(now),
		Lon:  pi.lonPred.GetEx(now),
		Alt:  pi.altPred.GetEx(now),
	}

	pi.pred.Set(pred)

	pi.pan.Tick()
	pi.tilt.Tick()
}

func (pi *PiPoint) predict(gps *Position) {
	now := pi.tick.GetFloat64()

	pi.latPred.SetEx(gps.Lat, now, gps.Time)
	pi.lonPred.SetEx(gps.Lon, now, gps.Time)
	pi.altPred.SetEx(gps.Alt, now, gps.Time)
}

func (p *PiPoint) update(param *Param) {
	state := p.state.GetInt()

	if state >= 0 && state < len(p.states) {
		p.states[state].Update(param)
	}

	p.log.Printf("%s %T %#v\n", param.Name, param.Get(), param.Get())
}

// Dispatch a MAVLink message.
func (p *PiPoint) Message(msg interface{}) {
	switch msg.(type) {
	case *common.Heartbeat:
		p.heartbeats.Inc()
		p.heartbeat.Set(msg.(*common.Heartbeat))
	case *common.SysStatus:
		p.sysStatus.Set(msg.(*common.SysStatus))
	case *common.GpsRawInt:
		gps := msg.(*common.GpsRawInt)
		p.gps.Set(&Position{
			Time:    float64(gps.TIME_USEC) * 1e-6,
			Lat:     float64(gps.LAT) * 1e-7,
			Lon:     float64(gps.LON) * 1e-7,
			Alt:     float64(gps.ALT) * 1e-3,
			Heading: float64(gps.COG) * 1e-2,
		})
		p.neu.Set(p.gps.Get().(*Position).ToNEU())
	case *common.Attitude:
		att := msg.(*common.Attitude)
		p.attitude.Set(&Attitude{
			float64(att.ROLL),
			float64(att.PITCH),
			float64(att.YAW),
		})
	default:
	}

	p.messages.Inc()
	p.log.Printf("%s %T %#v\n", "message", msg, msg)
}
