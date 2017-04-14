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

var (
	// Version is the overall binary version.  Set from the
	// build.
	Version = "dev"
)

// State is a handler for the current state of the system.
type State interface {
	Name() string
	Update(param *Param)
}

// PiPoint is an automatic, GPS based system that points a camera at the rover.
type PiPoint struct {
	Params *Params

	state *Param

	version    *Param
	tick       *Param
	seconds    *Param
	messages   *Param
	heartbeats *Param
	heartbeat  *Param
	attitude   *Param
	gps        *Param
	gpsFix     *Param
	neu        *Param
	baseOffset *Param
	pred       *Param
	rover      *Param
	base       *Param
	sysStatus  *Param
	link       *Param
	linkLast   int
	remote     *Param
	command    *Param
	mark       *Param
	vel        *Param

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

	audio *AudioOut
}

// NewPiPoint creates a new camera pointer.
func NewPiPoint() *PiPoint {
	p := &PiPoint{
		Params:  NewParams("pipoint"),
		latPred: &LinPred{},
		lonPred: &LinPred{},
		altPred: &LinPred{},
		elog:    NewEventLogger("pipoint"),
		param:   make(ParamChannel, 10),
		audio:   NewAudioOut(),
	}

	p.log = p.elog.logger

	p.states = []State{
		&LocateState{pi: p},
		&OrientateState{pi: p},
		&RunState{pi: p},
		&HoldState{pi: p},
		&CycleState{pi: p},
	}

	p.link = p.Params.NewNum("link.status")
	p.remote = p.Params.New("remote")
	p.command = p.Params.NewNum("command")
	p.mark = p.Params.NewNum("mark")

	p.version = p.Params.NewWith("build_label", Version)
	p.tick = p.Params.NewNum("tick")
	p.seconds = p.Params.NewNum("tick")
	p.messages = p.Params.NewNum("rover.messages")

	p.state = p.Params.NewNum("state")
	p.heartbeat = p.Params.NewWith("heartbeat", &common.Heartbeat{})
	p.heartbeats = p.Params.NewNum("heartbeat")

	p.gps = p.Params.New("gps")
	p.gpsFix = p.Params.NewNum("gps.fix")
	p.vel = p.Params.NewNum("gps.vog")
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

// AddMQTT adds a new MQTT connection that bridges between MQTT and
// params.
func (pi *PiPoint) AddMQTT(mqtt *mqtt.Adaptor) {
	NewParamMQTTBridge(pi.Params, mqtt, "")
}

// Run is the main entry point that runs forever.
func (pi *PiPoint) Run() {
	tick := time.NewTicker(dt)

	pi.audio.Say("Base ready")

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
	pi.seconds.UpdateInt(int(now))

	if !pi.heartbeat.Ok() {
		pi.link.UpdateInt(2)
	}

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

func (pi *PiPoint) update(param *Param) {
	if state := pi.getState(); state != nil {
		state.Update(param)
	}

	pi.announce(param)
	pi.log.Printf("%s %T %#v\n", param.Name, param.Get(), param.Get())
}

func (pi *PiPoint) getState() State {
	state := pi.state.GetInt()
	if state >= 0 && state < len(pi.states) {
		return pi.states[state]
	}
	return nil
}

func (pi *PiPoint) announce(param *Param) {
	switch param {
	case pi.state:
		if state := pi.getState(); state != nil {
			pi.audio.Say(state.Name())
		}
	case pi.link:
		link := param.GetInt()
		if link != pi.linkLast {
			pi.linkLast = link
			switch link {
			case 1:
				pi.audio.Say("Rover ready")
			case 2:
				pi.audio.Say("Rover offline")
			}
		}
	case pi.gpsFix:
		if param.GetInt() >= 3 {
			pi.audio.Say("GPS ready")
		}
	}
}

// Message handles a MAVLink message.
func (pi *PiPoint) Message(msg interface{}) {
	switch msg.(type) {
	case *common.Heartbeat:
		pi.link.SetInt(1)
		pi.heartbeats.Inc()
		pi.heartbeat.Set(msg.(*common.Heartbeat))
	case *common.SysStatus:
		pi.sysStatus.Set(msg.(*common.SysStatus))
	case *common.GpsRawInt:
		gps := msg.(*common.GpsRawInt)
		pi.gps.Set(&Position{
			Time:    float64(gps.TIME_USEC) * 1e-6,
			Lat:     float64(gps.LAT) * 1e-7,
			Lon:     float64(gps.LON) * 1e-7,
			Alt:     float64(gps.ALT) * 1e-3,
			Heading: float64(gps.COG) * 1e-2,
		})
		pi.neu.Set(pi.gps.Get().(*Position).ToNEU())
		pi.vel.SetFloat64(float64(gps.VEL) * 1e-2)
		pi.gpsFix.UpdateInt(int(gps.FIX_TYPE))
	case *common.Attitude:
		att := msg.(*common.Attitude)
		pi.attitude.Set(&Attitude{
			float64(att.ROLL),
			float64(att.PITCH),
			float64(att.YAW),
		})
	case *common.RcChannels:
		remote := msg.(*common.RcChannels)
		att := &Attitude{
			Roll:  ServoToScale(remote.CHAN1_RAW),
			Pitch: ServoToScale(remote.CHAN2_RAW),
			Yaw:   ServoToScale(remote.CHAN4_RAW),
		}
		pi.remote.Set(att)
		command := ScaleToPos(att.Yaw)
		if changed, _ := pi.command.UpdateInt(command); changed {
			if command == -2 {
				pi.mark.Inc()
			}
		}
	default:
	}

	pi.messages.Inc()
	pi.log.Printf("%s %T %#v\n", "message", msg, msg)
}
