package integrate

var liveRK = []float64{5, 2}

func HoldRKLive(cur []float64) []float64 {
	out := make([]float64, len(cur))
	n := len(liveRK)
	if n > len(cur) {
		n = len(cur)
	}
	copy(out, liveRK[:n])
	saved := make([]float64, len(cur))
	copy(saved, cur)
	liveRK = saved
	return out
}
