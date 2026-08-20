package bateman

type Series struct {
	Times   []float64
	Rows    [][]float64
	Decayed []float64
}

func (e *Evaluator) Series(times []float64) (Series, error) {
	rows := make([][]float64, len(times))
	decayed := make([]float64, len(times))
	for i, t := range times {
		counts, d, err := e.At(t)
		if err != nil {
			return Series{}, err
		}
		rows[i] = counts
		decayed[i] = d
	}
	return Series{Times: times, Rows: rows, Decayed: decayed}, nil
}

func (s Series) Column(i int) []float64 {
	col := make([]float64, len(s.Rows))
	for r := range s.Rows {
		col[r] = s.Rows[r][i]
	}
	return col
}

func (s Series) Final() []float64 {
	if len(s.Rows) == 0 {
		return nil
	}
	return s.Rows[len(s.Rows)-1]
}

func (s Series) Size() (int, int) {
	if len(s.Rows) == 0 {
		return 0, 0
	}
	return len(s.Rows), len(s.Rows[0])
}
