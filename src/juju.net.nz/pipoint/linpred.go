package pipoint

// A velocity based linear predictive filter.
type LinPred struct {
	x       float64
	stamp   float64
	updated float64
	v       float64
}

func (l *LinPred) Set(x float64) {
	now := Now()
	l.SetEx(x, now, now)
}

func (l *LinPred) SetEx(x, now, stamp float64) {
	if l.stamp == 0 {
		// First run
		l.v = 0
	} else {
		dt := stamp - l.stamp
		dx := x - l.x
		if dt <= 0 {
			l.v = 0
		} else {
			l.v = dx / dt
		}
	}
	l.x = x
	l.stamp = stamp
	l.updated = now
}

func (l *LinPred) Get() float64 {
	return l.GetEx(Now())
}

func (l *LinPred) GetEx(now float64) float64 {
	dt := now - l.updated
	if dt < 0 {
		dt = 0
	} else if dt > 2 {
		// Clamp if the value hasn't been updated recently.
		dt = 2
	}
	return l.x + l.v*dt
}
