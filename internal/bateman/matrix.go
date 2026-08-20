package bateman

func (s Series) Transpose() [][]float64 {
	if len(s.Rows) == 0 {
		return nil
	}
	n := len(s.Rows[0])
	out := make([][]float64, n)
	for i := 0; i < n; i++ {
		out[i] = make([]float64, len(s.Rows))
		for r := 0; r < len(s.Rows); r++ {
			out[i][r] = s.Rows[r][i]
		}
	}
	return out
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
