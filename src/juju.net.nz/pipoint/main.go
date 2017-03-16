package main

import (
	"fmt"
	"time"
	
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mavlink"
	common "gobot.io/x/gobot/platforms/mavlink/common"
)

func main() {
	adaptor := mavlink.NewAdaptor("192.168.0.93:14550")
	iris := mavlink.NewDriver(adaptor)

	work := func() {
		gobot.After(1*time.Second, func() {
			fmt.Println("after");

			dataStream := common.NewRequestDataStream(100,
				0,
				0,
				4,
				1,
			)
			iris.SendPacket(common.CraftMAVLinkPacket(0,
				0,
				dataStream,
			))
		})
		
		iris.Once(iris.Event(mavlink.PacketEvent), func(data interface{}) {
			fmt.Println("packet")
			packet := data.(*common.MAVLinkPacket)

			dataStream := common.NewRequestDataStream(100,
				packet.SystemID,
				packet.ComponentID,
				4,
				1,
			)
			iris.SendPacket(common.CraftMAVLinkPacket(packet.SystemID,
				packet.ComponentID,
				dataStream,
			))
		})

		iris.On(iris.Event(mavlink.MessageEvent), func(data interface{}) {
			fmt.Println("message")
			if data.(common.MAVLinkMessage).Id() == 30 {
				message := data.(*common.Attitude)
				fmt.Println("Attitude")
				fmt.Println("TIME_BOOT_MS", message.TIME_BOOT_MS)
				fmt.Println("ROLL", message.ROLL)
				fmt.Println("PITCH", message.PITCH)
				fmt.Println("YAW", message.YAW)
				fmt.Println("ROLLSPEED", message.ROLLSPEED)
				fmt.Println("PITCHSPEED", message.PITCHSPEED)
				fmt.Println("YAWSPEED", message.YAWSPEED)
				fmt.Println("")
			}
		})
	}

	robot := gobot.NewRobot("mavBot",
		[]gobot.Connection{adaptor},
		[]gobot.Device{iris},
		work,
	)

	fmt.Println("start")
	robot.Start()
}
