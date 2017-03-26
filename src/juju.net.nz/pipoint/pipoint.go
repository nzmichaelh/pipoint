package pipoint

import (
	"fmt"
	common "gobot.io/x/gobot/platforms/mavlink/common"
	"math"
	"log"
)

const (
	LocateState = iota
	RunState    = iota
	ManualState = iota
	CycleState  = iota
)

// Geographic coordinates.
type Position struct {
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

// An automatic, GPS based system that points a camera at the rover.
type PiPoint struct {
	Params *Params

	state *Param

	fix *Param

	heartbeats *Param
	heartbeat  *Param
	attitude   *Param
	gps        *Param
	rover      *Param
	base       *Param
	sysStatus  *Param

	sp     *Param
	offset *Param

	pan  *Servo
	tilt *Servo

	cycle float64
	// TODO: de-hack.
	lock  bool
}

// Create a new camera pointer.
func NewPiPoint() *PiPoint {
	p := &PiPoint{
		Params: NewParams("pipoint"),
	}

	p.state = p.Params.NewWith("state", 0)
	p.fix = p.Params.New("fix")
	p.heartbeat = p.Params.NewWith("heartbeat", &common.Heartbeat{})
	p.heartbeats = p.Params.New("heartbeat")

	p.gps = p.Params.New("gps.position")

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

func (p *PiPoint) Tick() {
	switch p.state.GetInt() {
	case CycleState:
		p.cycle += 0.02
		p.pan.Set(Scale(math.Cos(p.cycle), -1, 1, -math.Pi/2, math.Pi/2))
		p.tilt.Set(Scale(math.Sin(p.cycle), -1, 1, -math.Pi/2, 0))
	default:
	}

	p.pan.Tick()
	p.tilt.Tick()
}

func (p *PiPoint) check(code int, cond bool) bool {
	return !cond
}

func (p *PiPoint) update(param *Param) {
	switch p.state.GetInt() {
	case LocateState:
		p.locate(param)
	case RunState:
		p.run(param)
	}

	switch p.state.GetInt() {
	case LocateState, RunState:
		if param == p.sp || param == p.offset {
			sp := p.sp.Get().(*Attitude)
			offset := p.offset.Get().(*Attitude)

			p.pan.Set(sp.Yaw + offset.Yaw)
			p.tilt.Set(sp.Pitch + offset.Pitch)
		}
	default:
	}
}

// Updates during the Locate state.
func (p *PiPoint) locate(param *Param) {
	switch param {
	case p.gps:
		p.rover.Set(p.gps.Get())
		p.base.Set(p.gps.Get())
		p.base.Final()
	case p.attitude:
		p.offset.Set(&Attitude{
			Yaw: p.attitude.Get().(*Attitude).Yaw,
		})
	case p.fix:
		// Move to Run
		p.state.SetInt(RunState)
	}
}

func (p *PiPoint) point(rover, base *Position) (*Attitude, error) {
	lat := AsRad(base.Lat)

	dlat := rover.Lat - base.Lat
	dlon := rover.Lon - base.Lon
	dalt := rover.Alt - base.Alt

	if math.Abs(dlat) > 1 || math.Abs(dlon) > 1 {
		return nil, fmt.Errorf("Rover is too far away")
	}

	if math.Abs(lat) > AsRad(80) {
		return nil, fmt.Errorf("System is too far north or south")
	}

	dlat *= LatLength(lat)
	dlon *= LonLength(lat)

	hdist := math.Sqrt(dlat*dlat + dlon*dlon)
	pitch := math.Atan2(dalt, hdist)
	yaw := math.Atan2(dlon, dlat)

	return &Attitude{
		Pitch: pitch,
		Yaw:   yaw,
	}, nil
}

// Updates during the Run state.
func (p *PiPoint) run(param *Param) {
	switch param {
	case p.gps:
		p.rover.Set(p.gps.Get())
	}

	if !p.rover.Ok() || !p.base.Ok() {
		// Location is invalid or old.
		log.Println("run: skipping as invalid or old", p.rover.Ok(), p.base.Ok())
		return
	}

	rover := p.rover.Get().(*Position)
	base := p.base.Get().(*Position)

	att, err := p.point(rover, base)
	if err != nil {
		log.Printf("point: %v\n", err)
		return
	}

	if !p.lock {
		p.lock = true
		offset := p.offset.Get().(*Attitude)
		p.pan.Set(WrapAngle(att.Yaw + offset.Yaw))
		p.tilt.Set(WrapAngle(att.Pitch + offset.Pitch))
		p.lock = false
	}
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
			float64(gps.LAT) * 1e-7,
			float64(gps.LON) * 1e-7,
			float64(gps.ALT) * 1e-3,
			float64(gps.HDG) * 1e-2,
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

// Called when the user marks that the rover and base are at the same
// location.
func (p *PiPoint) Fix() {
	p.fix.Inc()
}
