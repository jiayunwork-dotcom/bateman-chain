package bateman

import "math"

func TwoNucleusParent(t, n0, parent float64) float64 {
	return n0 * math.Exp(-parent*t)
}

func TwoNucleusDaughter(t, n0, parent, daughter float64) float64 {
	if parent == daughter {
		return n0 * parent * t * math.Exp(-parent*t)
	}
	return n0 * parent * (math.Exp(-parent*t) - math.Exp(-daughter*t)) / (daughter - parent)
}

func EqualRateSeries(t, n0, rate float64, depth int) float64 {
	lt := rate * t
	p := 1.0
	for m := 1; m <= depth; m++ {
		p *= lt / float64(m)
	}
	return n0 * p * math.Exp(-rate*t)
}

func TwoNucleusIntegral(t, n0, parent, daughter float64) float64 {
	if parent == daughter {
		return n0 * (1 - (1+parent*t)*math.Exp(-parent*t))
	}
	den := daughter - parent
	return n0 * parent * ((1-math.Exp(-parent*t))/parent - (1-math.Exp(-daughter*t))/daughter) / den
}
