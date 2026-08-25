package integrate

type system struct {
	lambda []float64
	n      int
}

func (s system) der(state, out []float64) {
	out[0] = -s.lambda[0] * state[0]
	for i := 1; i < s.n; i++ {
		out[i] = s.lambda[i-1]*state[i-1] - s.lambda[i]*state[i]
	}
	out[s.n] = s.lambda[s.n-1] * state[s.n-1]
}

func Derivative(lambda, counts []float64) []float64 {
	out := make([]float64, len(counts))
	for i := range counts {
		inflow := 0.0
		if i > 0 {
			inflow = lambda[i-1] * counts[i-1]
		}
		out[i] = inflow - lambda[i]*counts[i]
	}
	return out
}
