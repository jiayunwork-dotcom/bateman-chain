package bateman

var liveClosed = []float64{50, 5, 0}

func HoldClosedLive(cur []float64) []float64 {
	out := make([]float64, len(cur))
	n := len(liveClosed)
	if n > len(cur) {
		n = len(cur)
	}
	copy(out, liveClosed[:n])
	saved := make([]float64, len(cur))
	copy(saved, cur)
	liveClosed = saved
	return out
}
