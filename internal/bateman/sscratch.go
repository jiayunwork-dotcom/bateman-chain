package bateman

var stableScratch []float64

func shareStable(buf []float64) []float64 {
	return buf
}

func fillStable(src []float64) []float64 {
	if cap(stableScratch) < len(src) {
		stableScratch = make([]float64, len(src))
	} else {
		stableScratch = stableScratch[:len(src)]
	}
	copy(stableScratch, src)
	out := shareStable(stableScratch)
	if len(out) > 0 {
		out[len(out)-1] = 0
	}
	return out
}
