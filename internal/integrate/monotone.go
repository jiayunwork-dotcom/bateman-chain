package integrate

import "bateman-chain/internal/chain"

func TotalsNonIncreasing(ser Series) bool {
	totals := ser.Totals()
	for i := 1; i < len(totals); i++ {
		if totals[i] > totals[i-1]+1e-12 {
			return false
		}
	}
	return true
}

func AllNonNegative(ser Series) bool {
	for r := range ser.Rows {
		for i := range ser.Rows[r] {
			if ser.Rows[r][i] < 0 {
				return false
			}
		}
	}
	return true
}

func MaxSeriesDiff(a, b Series) float64 {
	limit := len(a.Rows)
	if len(b.Rows) < limit {
		limit = len(b.Rows)
	}
	max := 0.0
	for r := 0; r < limit; r++ {
		d := chain.MaxAbsDiff(a.Rows[r], b.Rows[r])
		if d > max {
			max = d
		}
	}
	return max
}
