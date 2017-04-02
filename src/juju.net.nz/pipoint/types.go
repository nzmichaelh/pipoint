package pipoint

// Geographic coordinates.
type Position struct {
	Time    float64
	Lat     float64
	Lon     float64
	Alt     float64
	Heading float64
}

// On-surface coordinates.
type NEUPosition struct {
	Time  float64
	North float64
	East  float64
	Up   float64
}

// Orientation of a body.
type Attitude struct {
	Roll  float64
	Pitch float64
	Yaw   float64
}

func (p *Position) ToNEU() *NEUPosition {
	lat := AsRad(p.Lat)

	return &NEUPosition{
		Time: p.Time,
		North: LatLength(lat) * p.Lat,
		East: LonLength(lat) * p.Lon,
		Up: p.Alt,
	}
}

func (p *NEUPosition) Sub(right *NEUPosition) *NEUPosition {
	return &NEUPosition{
		Time: p.Time - right.Time,
		North: p.North - right.North,
		East: p.East - right.East,
		Up: p.Up - right.Up,
	}
}

func (p *NEUPosition) Add(right *NEUPosition) *NEUPosition {
	return &NEUPosition{
		Time: p.Time + right.Time,
		North: p.North + right.North,
		East: p.East + right.East,
		Up: p.Up + right.Up,
	}
}
