package pipoint

// A velocity based linear predictive filter.
type LinPred struct {
	x     float64
	stamp float64
	v     float64
}

func (l *LinPred) Set(x float64) {
	l.SetEx(x, Now())
}

func (l *LinPred) SetEx(x, now float64) {
	if l.stamp == 0 {
		// First run
		l.v = 0
	} else {
		dt := now - l.stamp
		dx := x - l.x
		if dt <= 0 {
			l.v = 0
		} else {
			l.v = dx / dt
		}
	}
	l.x = x
	l.stamp = now
}

func (l *LinPred) Get() float64 {
	return l.GetEx(Now())
}

func (l *LinPred) GetEx(now float64) float64 {
	dt := now - l.stamp
	if dt < 0 {
		dt = 0
	} else if dt > 2 {
		// Clamp if the value hasn't been updated recently.
		dt = 2
	}
	return l.x + l.v*dt
}
