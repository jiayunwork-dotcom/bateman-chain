package bateman

import (
	"math"

	"bateman-chain/internal/chain"
)

func VerifySeries(initial []float64, ser Series, relTol float64) error {
	initialTotal := chain.Total(initial)
	for i := range ser.Times {
		row := ser.Rows[i]
		if !chain.AllFinite(row) || math.IsNaN(ser.Decayed[i]) || math.IsInf(ser.Decayed[i], 0) {
			return chain.ErrNonFinite
		}
		remaining := chain.Total(row)
		if err := chain.CheckConservation(initialTotal, remaining, ser.Decayed[i], relTol); err != nil {
			return err
		}
	}
	return nil
}

func ClosestPair(lambda []float64) (int, int, float64) {
	bestI, bestJ := 0, 1
	best := math.Inf(1)
	for i := 0; i < len(lambda); i++ {
		for j := i + 1; j < len(lambda); j++ {
			gap := math.Abs(lambda[i] - lambda[j])
			if gap < best {
				best = gap
				bestI, bestJ = i, j
			}
		}
	}
	return bestI, bestJ, best
}
