package bateman

import "math"

const NearTol = 1e-8

func NearEqual(lambda []float64) bool {
	for i := 0; i < len(lambda); i++ {
		for j := i + 1; j < len(lambda); j++ {
			if nearly(lambda[i], lambda[j]) {
				return true
			}
		}
	}
	return false
}

func nearly(a, b float64) bool {
	scale := math.Max(math.Abs(a), math.Abs(b))
	if scale == 0 {
		return true
	}
	return math.Abs(a-b) <= NearTol*scale
}

func MinGap(lambda []float64) float64 {
	if len(lambda) < 2 {
		return 0
	}
	min := math.Inf(1)
	for i := 0; i < len(lambda); i++ {
		for j := i + 1; j < len(lambda); j++ {
			g := math.Abs(lambda[i] - lambda[j])
			if g < min {
				min = g
			}
		}
	}
	return min
}

func Separated(lambda []float64) bool {
	return !NearEqual(lambda)
}
