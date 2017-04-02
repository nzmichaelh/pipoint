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
	case s.pi.neu:
		s.pi.rover.Set(param.Get())
	}

	if !s.pi.rover.Ok() || !s.pi.base.Ok() {
		// Location is invalid or old.
		log.Println("run: skipping as invalid or old", s.pi.rover.Ok(), s.pi.base.Ok())
		return
	}

	rover := s.pi.rover.Get().(*NEUPosition)
	base := s.pi.base.Get().(*NEUPosition)
	offset := s.pi.baseOffset.Get().(*NEUPosition)

	att, err := point(rover, base, offset)
	if err != nil {
		log.Printf("point: %v\n", err)
		return
	}

	if param == s.pi.neu {
		offset := s.pi.offset.Get().(*Attitude)
		s.pi.pan.Set(WrapAngle(att.Yaw + offset.Yaw))
		s.pi.tilt.Set(WrapAngle(att.Pitch + offset.Pitch))
	}
}

func point(rover, base, offset *NEUPosition) (*Attitude, error) {
	delta := rover.Sub(base.Add(offset))
	if math.Abs(delta.North) > 10e3 || math.Abs(delta.East) > 10e3 {
		return nil, fmt.Errorf("Rover is too far away")
	}

	hdist := math.Sqrt(delta.North*delta.North + delta.East*delta.East)
	pitch := math.Atan2(delta.Up, hdist)
	yaw := math.Atan2(delta.East, delta.North)

	return &Attitude{
		Pitch: pitch,
		Yaw:   yaw,
	}, nil
}
