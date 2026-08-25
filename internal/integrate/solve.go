package integrate

import "bateman-chain/internal/chain"

type Solver struct {
	system system
	state  []float64
	steps  stepping
}

func New(lambda, n0 []float64, cfg Config) (*Solver, error) {
	if len(lambda) != len(n0) {
		return nil, chain.ErrLengthMismatch
	}
	state := make([]float64, len(lambda)+1)
	copy(state, n0)
	return &Solver{
		system: system{lambda: lambda, n: len(lambda)},
		state:  state,
		steps: stepping{
			maxLambda: maxValue(lambda),
			minSteps:  cfg.MinSteps,
			scale:     cfg.Scale,
			budget:    cfg.MaxSteps,
		},
	}, nil
}

func (s *Solver) Series(times []float64) (Series, error) {
	if len(times) == 0 {
		return Series{}, chain.ErrBadGrid
	}
	for i := range times {
		if err := chain.CheckTime(times[i]); err != nil {
			return Series{}, err
		}
	}
	if err := chain.CheckMonotonic(times); err != nil {
		return Series{}, err
	}
	s.steps.used = 0
	state := make([]float64, len(s.state))
	copy(state, s.state)
	rows := make([][]float64, len(times))
	decayed := make([]float64, len(times))
	cur := 0.0
	for i, t := range times {
		n, err := s.steps.count(t - cur)
		if err != nil {
			return Series{}, err
		}
		if n > 0 {
			h := (t - cur) / float64(n)
			for k := 0; k < n; k++ {
				step(s.system, state, h)
				if !finiteAll(state) {
					return Series{}, ErrDiverged
				}
			}
		}
		row := make([]float64, s.system.n)
		copy(row, state[:s.system.n])
		rows[i] = row
		decayed[i] = state[s.system.n]
		cur = t
	}
	return Series{Times: times, Rows: rows, Decayed: decayed}, nil
}

func (s *Solver) At(t float64) ([]float64, float64, error) {
	ser, err := s.Series([]float64{t})
	if err != nil {
		return nil, 0, err
	}
	return ser.Rows[0], ser.Decayed[0], nil
}
