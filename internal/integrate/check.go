package integrate

import "math"

func finiteAll(state []float64) bool {
	for i := 0; i < len(state); i++ {
		if math.IsNaN(state[i]) || math.IsInf(state[i], 0) {
			return false
		}
	}
	return true
}

func maxValue(lambda []float64) float64 {
	max := 0.0
	for i := 0; i < len(lambda); i++ {
		if lambda[i] > max {
			max = lambda[i]
		}
	}
	return max
}
