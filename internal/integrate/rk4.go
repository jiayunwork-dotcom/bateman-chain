package integrate

func step(s system, state []float64, h float64) {
	k1 := make([]float64, len(state))
	k2 := make([]float64, len(state))
	k3 := make([]float64, len(state))
	k4 := make([]float64, len(state))
	tmp := make([]float64, len(state))

	s.der(state, k1)
	for i := range state {
		tmp[i] = state[i] + 0.5*h*k1[i]
	}
	s.der(tmp, k2)
	for i := range state {
		tmp[i] = state[i] + 0.5*h*k2[i]
	}
	s.der(tmp, k3)
	for i := range state {
		tmp[i] = state[i] + h*k3[i]
	}
	s.der(tmp, k4)

	for i := range state {
		state[i] += h / 6 * (k1[i] + 2*k2[i] + 2*k3[i] + k4[i])
	}
}
