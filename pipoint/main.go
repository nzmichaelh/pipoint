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

package main

import (
	"log"
	"net/http"

	"juju.net.nz/x/pipoint"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/platforms/mavlink"
	"gobot.io/x/gobot/platforms/mqtt"
)

func main() {
	p := pipoint.NewPiPoint()

	mav := mavlink.NewUDPAdaptor(":14550")
	driver := mavlink.NewDriver(mav)

	// TODO(michaelh): make configurable.
	mq := mqtt.NewAdaptor("tls://iot.juju.net.nz:8883", "pipoint")
	mq.SetAutoReconnect(true)
	p.AddMQTT(mq)

	http.HandleFunc("/metrics", p.Params.Metrics)

	master := gobot.NewMaster()
	a := api.NewAPI(master)
	a.Start()

	robot := gobot.NewRobot("pipoint",
		[]gobot.Connection{mav, mq},
		[]gobot.Device{driver},
		func() {
			go p.Run()
			driver.On(driver.Event(mavlink.MessageEvent), p.Message)
		})

	master.AddRobot(robot)
	log.Println("Start")
	master.Start()
}