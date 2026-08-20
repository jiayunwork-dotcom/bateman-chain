package integrate

type Series struct {
	Times   []float64
	Rows    [][]float64
	Decayed []float64
}

func (s Series) Column(i int) []float64 {
	col := make([]float64, len(s.Rows))
	for r := range s.Rows {
		col[r] = s.Rows[r][i]
	}
	return col
}

func (s Series) Totals() []float64 {
	totals := make([]float64, len(s.Rows))
	for r := range s.Rows {
		var sum float64
		for i := 0; i < len(s.Rows[r]); i++ {
			sum += s.Rows[r][i]
		}
		totals[r] = sum
	}
	return totals
}
