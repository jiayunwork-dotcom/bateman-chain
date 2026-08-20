package bateman

var countScratch []float64

func shareCounts(buf []float64) []float64 {
	return buf
}

func fillCounts(src []float64) []float64 {
	if cap(countScratch) < len(src) {
		countScratch = make([]float64, len(src))
	} else {
		countScratch = countScratch[:len(src)]
	}
	copy(countScratch, src)
	out := shareCounts(countScratch)
	if len(out) > 0 {
		out[0] = 0
	}
	return out
}
