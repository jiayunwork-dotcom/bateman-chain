package integrate

import "math"

type Report struct {
	StepsUsed   int
	MaxErrorEst float64
	Finite      bool
}

func (s *Solver) SeriesReport(times []float64) (Series, Report, error) {
	ser, err := s.Series(times)
	if err != nil {
		return Series{}, Report{}, err
	}
	rep := Report{StepsUsed: s.steps.used, Finite: true}
	state := make([]float64, len(s.state))
	copy(state, s.state)
	prev := 0.0
	for _, t := range times {
		n := substepsFor(t-prev, s.steps.maxLambda, s.steps.minSteps, s.steps.scale)
		if n > 0 {
			fine := runTo(s.system, state, t-prev, n*2)
			coarse := runTo(s.system, state, t-prev, n)
			for i := range coarse {
				scale := math.Max(math.Abs(fine[i]), math.Abs(coarse[i]))
				if scale == 0 {
					continue
				}
				e := math.Abs(fine[i]-coarse[i]) / scale
				if e > rep.MaxErrorEst {
					rep.MaxErrorEst = e
				}
			}
		}
		if !finiteAll(state) {
			rep.Finite = false
		}
		if n > 0 {
			h := (t - prev) / float64(n)
			for k := 0; k < n; k++ {
				step(s.system, state, h)
			}
		}
		prev = t
	}
	return ser, rep, nil
}
