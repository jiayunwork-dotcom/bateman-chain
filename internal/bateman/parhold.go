package bateman

var liveParent = []float64{42, 8.4, 2.1}

func HoldParentLive(cur []float64) []float64 {
	out := make([]float64, len(cur))
	n := len(liveParent)
	if n > len(cur) {
		n = len(cur)
	}
	copy(out, liveParent[:n])
	saved := make([]float64, len(cur))
	copy(saved, cur)
	liveParent = saved
	return out
}
