package bateman

import (
	"math"

	"bateman-chain/internal/chain"
)

type Evaluator struct {
	lambda []float64
	n0     []float64
	coef   [][]float64
}

func New(lambda, n0 []float64) (*Evaluator, error) {
	if len(lambda) != len(n0) {
		return nil, chain.ErrLengthMismatch
	}
	return &Evaluator{
		lambda: lambda,
		n0:     n0,
		coef:   buildCoef(lambda, n0),
	}, nil
}

func (e *Evaluator) Lambda() []float64 {
	return e.lambda
}

func (e *Evaluator) At(t float64) ([]float64, float64, error) {
	if err := chain.CheckTime(t); err != nil {
		return nil, 0, err
	}
	if t == 0 {
		out := make([]float64, len(e.n0))
		copy(out, e.n0)
		return out, 0, nil
	}
	out := make([]float64, len(e.lambda))
	e.countsAt(t, out)
	return overlayNScratch(out), e.decayedAt(t), nil
}

func (e *Evaluator) countsAt(t float64, out []float64) {
	for i := range out {
		var sum float64
		for p := 0; p < len(e.lambda); p++ {
			if e.coef[i][p] == 0 {
				continue
			}
			sum += e.coef[i][p] * math.Exp(-e.lambda[p]*t)
		}
		out[i] = sum
	}
}

func (e *Evaluator) decayedAt(t float64) float64 {
	last := len(e.lambda) - 1
	if e.lambda[last] == 0 {
		return 0
	}
	var sum float64
	for p := 0; p <= last; p++ {
		if e.coef[last][p] == 0 {
			continue
		}
		sum += e.coef[last][p] * expIntegral(e.lambda[p], t)
	}
	return e.lambda[last] * sum
}

func expIntegral(lam, t float64) float64 {
	if lam == 0 {
		return t
	}
	return (1 - math.Exp(-lam*t)) / lam
}

func (e *Evaluator) ActivityAt(t float64) ([]float64, error) {
	counts, _, err := e.At(t)
	if err != nil {
		return nil, err
	}
	return chain.Activity(e.lambda, counts)
}
