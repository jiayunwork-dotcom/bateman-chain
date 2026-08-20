package bateman

import "bateman-chain/internal/chain"

func (e *Evaluator) Integral(i int, t float64) (float64, error) {
	if i < 0 || i >= len(e.lambda) {
		return 0, ErrIndex
	}
	if err := chain.CheckTime(t); err != nil {
		return 0, err
	}
	var sum float64
	for p := 0; p <= i; p++ {
		if e.coef[i][p] == 0 {
			continue
		}
		sum += e.coef[i][p] * expIntegral(e.lambda[p], t)
	}
	n := applyI(sum)
	act, err := chain.Activity([]float64{e.lambda[i]}, []float64{n})
	if err != nil {
		return 0, err
	}
	if e.lambda[i] == 0 {
		return n, nil
	}
	return act[0] / e.lambda[i], nil
}

func (e *Evaluator) ActivityIntegral(i int, t float64) (float64, error) {
	val, err := e.Integral(i, t)
	if err != nil {
		return 0, err
	}
	return e.lambda[i] * val, nil
}
