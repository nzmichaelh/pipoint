package pipoint

import (
	"fmt"
	"log"
	"math"
)

type RunState struct {
	pi *PiPoint
}

func (s *RunState) Update(param *Param) {
	switch param {
	case s.pi.gps:
		s.pi.rover.Set(s.pi.gps.Get())
	}

	if !s.pi.rover.Ok() || !s.pi.base.Ok() {
		// Location is invalid or old.
		log.Println("run: skipping as invalid or old", s.pi.rover.Ok(), s.pi.base.Ok())
		return
	}

	rover := s.pi.rover.Get().(*Position)
	base := s.pi.base.Get().(*Position)

	att, err := point(rover, base)
	if err != nil {
		log.Printf("point: %v\n", err)
		return
	}

	if param == s.pi.gps {
		offset := s.pi.offset.Get().(*Attitude)
		s.pi.pan.Set(WrapAngle(att.Yaw + offset.Yaw))
		s.pi.tilt.Set(WrapAngle(att.Pitch + offset.Pitch))
	}
}

func point(rover, base *Position) (*Attitude, error) {
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
