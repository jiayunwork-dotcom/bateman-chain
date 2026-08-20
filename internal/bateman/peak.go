package bateman

import (
	"math"

	"bateman-chain/internal/chain"
)

func FirstDaughterPeakTime(parent, daughter float64) float64 {
	if parent == daughter {
		return 1 / parent
	}
	return math.Log(parent/daughter) / (parent - daughter)
}

func FirstDaughterPeakValue(initial, parent, daughter float64) float64 {
	if parent == daughter {
		return initial * math.Exp(-1)
	}
	peak := FirstDaughterPeakTime(parent, daughter)
	return initial * parent * (math.Exp(-parent*peak) - math.Exp(-daughter*peak)) / (daughter - parent)
}

func ScanPeak(e *Evaluator, times []float64) (int, float64) {
	best := 0
	var bestVal float64
	for i := 0; i < len(times); i++ {
		row, _, err := e.At(times[i])
		if err != nil {
			return -1, math.NaN()
		}
		val := row[len(row)-1]
		if i == 0 || val > bestVal {
			best = i
			bestVal = val
		}
	}
	return best, bestVal
}

func BuildupTime(s chain.Spec, fraction float64) float64 {
	longest := 0.0
	for i := 0; i < s.Len(); i++ {
		if s.Lambda[i] == 0 {
			continue
		}
		lt := chain.MeanLifetime(s.Lambda[i])
		if lt > longest {
			longest = lt
		}
	}
	if longest == 0 {
		return 0
	}
	return -longest * math.Log(1-fraction)
}
