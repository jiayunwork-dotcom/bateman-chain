package bateman

func denominator(lambda []float64, from, target, skip int) float64 {
	den := 1.0
	for q := from; q <= target; q++ {
		if q == skip {
			continue
		}
		den *= lambda[q] - lambda[skip]
	}
	return den
}
