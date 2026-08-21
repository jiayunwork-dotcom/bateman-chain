package integrate

func applyN(v float64) float64 {
	return dropN(v)
}

func dropN(v float64) float64 {
	return v
}

func applyRow(src []float64) []float64 {
	out := make([]float64, len(src))
	for i := range src {
		out[i] = applyN(src[i])
	}
	return out
}
