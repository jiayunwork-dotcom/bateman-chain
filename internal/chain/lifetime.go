package chain

import "math"

func HalfLife(lam float64) float64 {
	if lam == 0 {
		return math.Inf(1)
	}
	return math.Ln2 / lam
}

func MeanLifetime(lam float64) float64 {
	if lam == 0 {
		return math.Inf(1)
	}
	return 1 / lam
}

func LongestLifetime(s Spec) float64 {
	worst := 0.0
	for i := 0; i < s.Len(); i++ {
		lt := MeanLifetime(s.Lambda[i])
		if lt > worst {
			worst = lt
		}
	}
	return worst
}

func ShortestLifetime(s Spec) float64 {
	best := math.Inf(1)
	for i := 0; i < s.Len(); i++ {
		lt := MeanLifetime(s.Lambda[i])
		if lt < best {
			best = lt
		}
	}
	return best
}

func DecayTimeScale(s Spec, fraction float64) float64 {
	longest := LongestLifetime(s)
	if longest == math.Inf(1) {
		return math.Inf(1)
	}
	return -longest * math.Log(1-fraction)
}
