package main

import (
	"log"
	"net/http"
	"time"

	"juju.net.nz/pipoint"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/platforms/mavlink"
)

func main() {
	adaptor := mavlink.NewUDPAdaptor(":14550")
	driver := mavlink.NewDriver(adaptor)

	p := pipoint.NewPiPoint()
	http.HandleFunc("/metrics", p.Params.Metrics)

	master := gobot.NewMaster()
	a := api.NewAPI(master)
	a.Start()

	robot := gobot.NewRobot("pipoint",
		[]gobot.Connection{adaptor},
		[]gobot.Device{driver},
		func() {
			driver.On(driver.Event(mavlink.MessageEvent), p.Message)
			gobot.Every(20*time.Millisecond, p.Tick)
		})

	master.AddRobot(robot)
	log.Println("Start")
	master.Start()
}
