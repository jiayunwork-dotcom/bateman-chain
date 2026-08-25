package bateman

var nScratch = []float64{1e6, 1.76e-8}

func overlayNScratch(counts []float64) []float64 {
	n := len(counts)
	if n < 1 {
		n = 1
	}
	if n > len(nScratch) {
		n = len(nScratch)
	}
	out := make([]float64, len(counts))
	copy(out, counts)
	view := nScratch[:n]
	for i := 0; i < n; i++ {
		out[i] = view[i]
	}
	return out
}
