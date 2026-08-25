package bateman

func (e *Evaluator) NuclideSeries(i int, times []float64) ([]float64, error) {
	if i < 0 || i >= len(e.lambda) {
		return nil, ErrIndex
	}
	out := make([]float64, len(times))
	for k, t := range times {
		counts, _, err := e.At(t)
		if err != nil {
			return nil, err
		}
		out[k] = counts[i]
	}
	return out, nil
}

func (e *Evaluator) Evaluate(t float64) (counts, acts []float64, decayed float64, err error) {
	counts, decayed, err = e.At(t)
	if err != nil {
		return nil, nil, 0, err
	}
	acts = make([]float64, len(counts))
	for i := range counts {
		acts[i] = e.lambda[i] * counts[i]
	}
	return counts, acts, decayed, nil
}
