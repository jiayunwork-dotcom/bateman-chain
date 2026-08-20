package integrate

import (
	"math"

	"bateman-chain/internal/chain"
)

type Accuracy struct {
	BaseSteps int
	RefinedSteps int
	MaxError float64
}

func EstimateFinalError(sys system, state []float64, t, scale float64, base int) Accuracy {
	refined := base * 2
	coarse := runTo(sys, state, t, base)
	fine := runTo(sys, state, t, refined)
	err := 0.0
	for i := range coarse {
		den := math.Max(math.Abs(fine[i]), scale)
		if den == 0 {
			continue
		}
		e := math.Abs(coarse[i]-fine[i]) / den
		if e > err {
			err = e
		}
	}
	return Accuracy{BaseSteps: base, RefinedSteps: refined, MaxError: err}
}

func runTo(sys system, state []float64, t float64, n int) []float64 {
	work := make([]float64, len(state))
	copy(work, state)
	h := t / float64(n)
	for i := 0; i < n; i++ {
		step(sys, work, h)
	}
	return work
}

func StepsForTolerance(sys system, state []float64, t, scale, tol float64, base, budget int) (int, error) {
	n := base
	for n*2 <= budget {
		acc := EstimateFinalError(sys, state, t, scale, n)
		if acc.MaxError <= tol {
			return n, nil
		}
		n *= 2
	}
	return 0, chain.ErrStepLimit
}
