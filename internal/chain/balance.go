package chain

import "math"

type Conservation struct {
	Initial   float64
	Remaining float64
	Decayed   float64
}

func Total(counts []float64) float64 {
	var sum float64
	for i := 0; i < len(counts); i++ {
		sum += counts[i]
	}
	return sum
}

func InitialTotal(s Spec) float64 {
	return Total(s.N0)
}

func (c Conservation) Residual() float64 {
	return c.Initial - c.Remaining - c.Decayed
}

func CheckConservation(initial, remaining, decayed, relTol float64) error {
	scale := math.Abs(initial)
	if scale < 1 {
		scale = 1
	}
	cons := Conservation{Initial: initial, Remaining: remaining, Decayed: decayed}
	if math.Abs(cons.Residual()) > relTol*scale {
		return ErrConservation
	}
	return nil
}
