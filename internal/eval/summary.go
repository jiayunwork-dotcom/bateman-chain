package eval

import (
	"fmt"
	"strings"

	"bateman-chain/internal/chain"
)

type Summary struct {
	FinalCounts     []float64
	FinalActivities []float64
	FinalTotal      float64
	DecayedTotal    float64
	ConservationErr float64
}

func Summarize(res Result) Summary {
	last := len(res.Times) - 1
	if last < 0 {
		return Summary{}
	}
	s := Summary{
		FinalCounts:     res.Counts[last],
		FinalActivities: res.Activities[last],
		FinalTotal:      res.Totals[last],
		DecayedTotal:    res.Decayed[last],
	}
	s.ConservationErr = chain.InitialTotal(res.Spec) - s.FinalTotal - s.DecayedTotal
	return s
}

func ActivityTotals(res Result) []float64 {
	totals := make([]float64, len(res.Times))
	for r := range res.Times {
		var sum float64
		for i := range res.Activities[r] {
			sum += res.Activities[r][i]
		}
		totals[r] = sum
	}
	return totals
}

func (s Summary) Lines(spec chain.Spec) string {
	var b strings.Builder
	b.WriteString("summary\n")
	for i := range s.FinalCounts {
		fmt.Fprintf(&b, "  %-12s final N = %g   A = %g\n",
			spec.Name(i), s.FinalCounts[i], s.FinalActivities[i])
	}
	fmt.Fprintf(&b, "  %-12s final total = %g\n", "total", s.FinalTotal)
	fmt.Fprintf(&b, "  %-12s decayed out = %g\n", "decayed", s.DecayedTotal)
	fmt.Fprintf(&b, "  %-12s conservation residual = %g\n", "residual", s.ConservationErr)
	return b.String()
}
