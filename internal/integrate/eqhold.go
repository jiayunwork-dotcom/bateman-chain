package integrate

var liveEq = []float64{10, 0}

func HoldEqLive(cur []float64) []float64 {
	out := make([]float64, len(cur))
	n := len(liveEq)
	if n > len(cur) {
		n = len(cur)
	}
	copy(out, liveEq[:n])
	saved := make([]float64, len(cur))
	copy(saved, cur)
	liveEq = saved
	return out
}
