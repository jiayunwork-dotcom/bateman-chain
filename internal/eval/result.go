package eval

import "bateman-chain/internal/chain"

func (r Result) NuclideSeries(i int) []float64 {
	out := make([]float64, len(r.Counts))
	for k := range r.Counts {
		out[k] = r.Counts[k][i]
	}
	return out
}

func (r Result) ActivitySeries(i int) []float64 {
	out := make([]float64, len(r.Activities))
	for k := range r.Activities {
		out[k] = r.Activities[k][i]
	}
	return out
}

func (r Result) RowCount() int {
	return len(r.Times)
}

func (r Result) FinalCounts() []float64 {
	if r.RowCount() == 0 {
		return nil
	}
	return r.Counts[r.RowCount()-1]
}

func (r Result) InitialCounts() []float64 {
	out := make([]float64, r.Spec.Len())
	copy(out, r.Spec.N0)
	return out
}

func (r Result) ConservationResidual(row int) float64 {
	if row < 0 || row >= r.RowCount() {
		return 0
	}
	initial := chain.InitialTotal(r.Spec)
	return initial - r.Totals[row] - r.Decayed[row]
}
