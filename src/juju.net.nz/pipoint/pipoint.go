package pipoint

import (
	common "gobot.io/x/gobot/platforms/mavlink/common"
)

const (
	LocateState = iota
	RunState    = iota
)

// Geographic coordinates.
type Position struct {
	Lat float64
	Lon float64
	Alt float64
}

// Orientation of a body.
type Attitude struct {
	Roll  float64
	Pitch float64
	Yaw   float64
}

// An automatic, GPS based system that points a camera at the rover.
type PiPoint struct {
	params *Params

	state *Param

	fix *Param

	heartbeat *Param
	attitude  *Param
	gps       *Param
	rover     *Param
	base      *Param

	sp     *Param
	offset *Param

	pan  *Param
	tilt *Param
}

// Create a new camera pointer.
func NewPiPoint() *PiPoint {
	p := &PiPoint{
		params: &Params{},
	}

	p.state = p.params.NewWith("state", 0)
	p.fix = p.params.New("fix")
	p.heartbeat = p.params.New("heartbeat")

	p.gps = p.params.New("gps.position")

	p.attitude = p.params.New("rover.attitude")
	p.rover = p.params.New("rover.position")
	p.base = p.params.New("base.position")

	p.sp = p.params.NewWith("pantilt.sp", &Attitude{})
	p.offset = p.params.NewWith("pantilt.offset", &Attitude{})

	p.pan = p.params.New("pantilt.pan")
	p.tilt = p.params.New("pantilt.tilt")

	p.params.Listen(p.update)
	return p
}

func (p *PiPoint) update(param *Param) {
	switch p.state.GetInt() {
	case LocateState:
		p.locate(param)
	case RunState:
		p.run(param)
	}

	if param == p.sp || param == p.offset {
		sp := p.sp.Get().(*Attitude)
		offset := p.offset.Get().(*Attitude)

		p.pan.SetFloat64(sp.Yaw + offset.Yaw)
		p.tilt.SetFloat64(sp.Pitch + offset.Pitch)
	}
}

// Updates during the Locate state.
func (p *PiPoint) locate(param *Param) {
	switch param {
	case p.gps:
		p.rover.Set(p.gps.Get())
	case p.attitude:
		p.offset.Set(&Attitude{
			Yaw: p.attitude.Get().(*Attitude).Yaw,
		})
	case p.fix:
		// Move to Run
		p.state.SetInt(RunState)
	}
}

// Updates during the Run state.
func (p *PiPoint) run(param *Param) {
}

// Dispatch a MAVLink message.
func (p *PiPoint) Message(msg interface{}) {
	switch msg.(type) {
	case *common.Heartbeat:
		p.heartbeat.Inc()
	case *common.GpsRawInt:
		gps := msg.(*common.GpsRawInt)
		p.gps.Set(&Position{
			float64(gps.LAT),
			float64(gps.LON),
			float64(gps.ALT),
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
