package eval

import (
	"fmt"

	"bateman-chain/internal/bateman"
	"bateman-chain/internal/chain"
)

type CrossCheck struct {
	Times      []float64
	Closed     [][]float64
	Numeric    [][]float64
	MaxRelDiff float64
	WithinTol  bool
}

func CrossCheckClosedNumeric(c Case, times []float64, tol float64, maxSteps int) (CrossCheck, error) {
	req := Request{Times: times, StepsLimit: chain.GridStepsLimit}
	closed, err := RunOpts(c, req, Options{Solver: SolverBateman})
	if err != nil {
		return CrossCheck{}, err
	}
	numeric, err := RunOpts(c, req, Options{Solver: SolverRK4, IntegrateMaxSteps: maxSteps})
	if err != nil {
		return CrossCheck{}, err
	}
	for r := range closed.Counts {
		closed.Counts[r] = bateman.HoldClosedLive(closed.Counts[r])
	}
	cc := CrossCheck{Times: times, Closed: closed.Counts, Numeric: numeric.Counts}
	for r := range times {
		rel := chain.RelativeDiff(closed.Counts[r], numeric.Counts[r])
		if rel > cc.MaxRelDiff {
			cc.MaxRelDiff = rel
		}
	}
	cc.WithinTol = cc.MaxRelDiff <= tol
	return cc, nil
}

func (cc CrossCheck) Describe() string {
	return fmt.Sprintf("max relative diff %.3e (within tolerance %v)", cc.MaxRelDiff, cc.WithinTol)
}
