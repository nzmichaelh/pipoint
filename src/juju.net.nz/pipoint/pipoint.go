package pipoint

import (
	"log"
	
	common "gobot.io/x/gobot/platforms/mavlink/common"
)

// Geographic coordinates.
type Position struct {
	Time    float64
	Lat     float64
	Lon     float64
	Alt     float64
	Heading float64
}

// Orientation of a body.
type Attitude struct {
	Roll  float64
	Pitch float64
	Yaw   float64
}

type State interface {
	Update(param *Param)
}

// An automatic, GPS based system that points a camera at the rover.
type PiPoint struct {
	Params *Params

	state *Param

	tick       *Param
	heartbeats *Param
	heartbeat  *Param
	attitude   *Param
	gps        *Param
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

	cycle  float64
	states []State

	elog *EventLogger
	log *log.Logger
}

// Create a new camera pointer.
func NewPiPoint() *PiPoint {
	p := &PiPoint{
		Params: NewParams("pipoint"),
		latPred: &LinPred{},
		lonPred: &LinPred{},
		altPred: &LinPred{},
		elog: NewEventLogger("pipoint"),
	}

	p.log = p.elog.logger
	
	p.states = []State{
		&LocateState{pi: p},
		&RunState{pi: p},
		&CycleState{pi: p},
	}

	p.tick = p.Params.New("tick")

	p.state = p.Params.NewWith("state", 0)
	p.heartbeat = p.Params.NewWith("heartbeat", &common.Heartbeat{})
	p.heartbeats = p.Params.New("heartbeat")

	p.gps = p.Params.New("gps.position")
	p.pred = p.Params.New("pred.position")

	p.attitude = p.Params.New("rover.attitude")
	p.rover = p.Params.New("rover.position")
	p.base = p.Params.New("base.position")

	p.sysStatus = p.Params.New("rover.status")

	p.sp = p.Params.NewWith("pantilt.sp", &Attitude{})
	p.offset = p.Params.NewWith("pantilt.offset", &Attitude{})

	p.pan = NewServo("pantilt.pan", p.Params)
	p.tilt = NewServo("pantilt.tilt", p.Params)

	p.Params.Listen(p.update)
	p.Params.Load()
	return p
}

func (pi *PiPoint) Tick() {
	now := Now()
	pi.tick.SetFloat64(now)

	pred := &Position{
		Time: now,
		Lat: pi.latPred.GetEx(now),
		Lon: pi.lonPred.GetEx(now),
		Alt: pi.altPred.GetEx(now),
	}

	pi.pred.Set(pred)
	
	pi.pan.Tick()
	pi.tilt.Tick()
}

func (p *PiPoint) check(code int, cond bool) bool {
	return !cond
}

func (pi *PiPoint) predict(gps *Position) {
	now := Now()

	pi.latPred.SetEx(gps.Lat, now)
	pi.lonPred.SetEx(gps.Lon, now)
	pi.altPred.SetEx(gps.Alt, now)
}

func (p *PiPoint) update(param *Param) {
	switch param {
	case p.gps:
		if param.Ok() {
			p.predict(param.Get().(*Position))
		}
	}

	state := p.state.GetInt()

	if state >= 0 && state < len(p.states) {
		p.states[state].Update(param)
	}

	p.log.Printf("%s %T %#v\n", param.name, param.Get(), param.Get())
}

// Dispatch a MAVLink message.
func (p *PiPoint) Message(msg interface{}) {
	switch msg.(type) {
	case *common.Heartbeat:
		p.heartbeats.Inc()
		p.heartbeat.Set(msg.(*common.Heartbeat))
	case *common.SysStatus:
		p.sysStatus.Set(msg.(*common.SysStatus))
	case *common.GlobalPositionInt:
		gps := msg.(*common.GlobalPositionInt)
		p.gps.Set(&Position{
			Time:    float64(gps.TIME_BOOT_MS) * 1e-3,
			Lat:     float64(gps.LAT) * 1e-7,
			Lon:     float64(gps.LON) * 1e-7,
			Alt:     float64(gps.ALT) * 1e-3,
			Heading: float64(gps.HDG) * 1e-2,
		})
	case *common.Attitude:
		att := msg.(*common.Attitude)
		p.attitude.Set(&Attitude{
			float64(att.ROLL),
			float64(att.PITCH),
			float64(att.YAW),
		})
	default:
	}
}
