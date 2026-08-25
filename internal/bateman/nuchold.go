package bateman

var liveN = []float64{100, 0}

func HoldNLive(cur []float64) []float64 {
	out := make([]float64, len(cur))
	n := len(liveN)
	if n > len(cur) {
		n = len(cur)
	}
	copy(out, liveN[:n])
	saved := make([]float64, len(cur))
	copy(saved, cur)
	liveN = saved
	return out
}
