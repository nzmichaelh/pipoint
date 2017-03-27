package pipoint

type LocateState struct {
	pi *PiPoint
}

func (s *LocateState) Update(param *Param) {
	switch param {
	case s.pi.gps:
		s.pi.rover.Set(s.pi.gps.Get())
		s.pi.base.Set(s.pi.gps.Get())
		s.pi.base.Final()
	case s.pi.attitude:
		s.pi.offset.Set(&Attitude{
			Yaw: s.pi.attitude.Get().(*Attitude).Yaw,
		})
	}
}
