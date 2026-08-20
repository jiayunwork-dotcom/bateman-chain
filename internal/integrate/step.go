package integrate

type StepSizes struct {
	Intervals []int
	Substep   float64
	Total     int
}

func (s *Solver) StepSizes(times []float64) StepSizes {
	cfg := s.steps
	prev := 0.0
	ss := StepSizes{}
	ss.Intervals = make([]int, len(times))
	for i, t := range times {
		dt := t - prev
		n := substepsFor(dt, cfg.maxLambda, cfg.minSteps, cfg.scale)
		ss.Intervals[i] = n
		ss.Total += n
		if n > 0 {
			ss.Substep = dt / float64(n)
		}
		prev = t
	}
	return ss
}
