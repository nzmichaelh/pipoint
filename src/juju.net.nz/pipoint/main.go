package main

import (
	"log"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	_ "gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/mavlink"
	common "gobot.io/x/gobot/platforms/mavlink/common"
)

const (
	LocateState = iota
	RunState    = iota
)

type Position struct {
	Lat float64
	Lon float64
	Alt float64
}

type Attitude struct {
	Roll  float64
	Pitch float64
	Yaw   float64
}

type PiPoint struct {
	params *Params
	driver *mavlink.Driver

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

func NewPiPoint(driver *mavlink.Driver) *PiPoint {
	p := &PiPoint{
		params: &Params{},
		driver: driver,
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

	p.params.Listen(p.Update)
	return p
}

func (p *PiPoint) Update(param *Param) {
	switch p.state.GetInt() {
	case LocateState:
		p.Locate(param)
	case RunState:
		p.Run(param)
	}

	if param == p.sp || param == p.offset {
		sp := p.sp.Get().(*Attitude)
		offset := p.offset.Get().(*Attitude)

		p.pan.SetFloat64(sp.Yaw + offset.Yaw)
		p.tilt.SetFloat64(sp.Pitch + offset.Pitch)
	}
}

func (p *PiPoint) Locate(param *Param) {
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

func (p *PiPoint) Run(param *Param) {
}

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

func (p *PiPoint) Fix() string {
	p.fix.Inc()
	return "OK"
}

func (p *PiPoint) Work() {
	p.driver.On(p.driver.Event(mavlink.MessageEvent), p.Message)
}

func main() {
	adaptor := mavlink.NewUDPAdaptor(":14550")
	driver := mavlink.NewDriver(adaptor)

	p := NewPiPoint(driver)

	master := gobot.NewMaster()
	a := api.NewAPI(master)
	a.Start()

	robot := gobot.NewRobot("pipoint",
		[]gobot.Connection{adaptor},
		[]gobot.Device{driver},
		func() { p.Work() },
	)

	robot.AddCommand("fix", func(params map[string]interface{}) interface{} {
		return p.Fix()
	})

	master.AddRobot(robot)
	log.Println("Start")
	master.Start()
}
