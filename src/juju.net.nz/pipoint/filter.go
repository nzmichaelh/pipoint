package pipoint

type Lowpass struct {
	Tau float64
	Acc float64
}

func (l *Lowpass) Step(v float64) float64 {
	return l.StepEx(v, l.Tau)
}

func (l *Lowpass) StepEx(v, tau float64) float64 {
	l.Acc = tau*v + (1-tau)*l.Acc
	return l.Acc
}
