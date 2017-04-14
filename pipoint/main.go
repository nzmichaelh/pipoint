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
	"flag"
	"net/http"

	"juju.net.nz/x/pipoint"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/platforms/mavlink"
	"gobot.io/x/gobot/platforms/mqtt"
)

func main() {
	mqttUrl := flag.String("mqtt.url", "", "URI of the MQTT server, such as tls://iot.juju.net.nz:8883")
	mavAddr := flag.String("mavlink.address", ":14550", "Address to listen on for Mavlink messages")

	flag.Parse()

	var cons []gobot.Connection
	var drivers []gobot.Device

	pi := pipoint.NewPiPoint()

	if mavAddr != nil && *mavAddr != "" {
		mav := mavlink.NewUDPAdaptor(*mavAddr)
		cons = append(cons, mav)
		driver := mavlink.NewDriver(mav)
		drivers = append(drivers, driver)
		driver.On(driver.Event(mavlink.MessageEvent), pi.Message)
	}

	if mqttUrl != nil && *mqttUrl != "" {
		mq := mqtt.NewAdaptor(*mqttUrl, "pipoint")
		cons = append(cons, mq)
		mq.SetAutoReconnect(true)
		pi.AddMQTT(mq)
	}

	http.HandleFunc("/metrics", pi.Params.Metrics)

	master := gobot.NewMaster()
	a := api.NewAPI(master)
	a.Start()

	robot := gobot.NewRobot("pipoint",
		cons, drivers,
		func() {
			go pi.Run()
		})

	master.AddRobot(robot)
	master.Start()
}
