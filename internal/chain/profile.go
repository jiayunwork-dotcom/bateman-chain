package chain

type Profile struct {
	Rates   []float64
	Counts  []float64
	Total   float64
}

func (s Spec) Profile() Profile {
	total := InitialTotal(s)
	rates := make([]float64, s.Len())
	counts := make([]float64, s.Len())
	for i := 0; i < s.Len(); i++ {
		rates[i] = s.Lambda[i]
		counts[i] = s.N0[i]
	}
	return Profile{Rates: rates, Counts: counts, Total: total}
}

func (p Profile) ParentShare() float64 {
	if p.Total == 0 {
		return 0
	}
	return p.Counts[0] / p.Total
}

func (p Profile) TerminalShare() float64 {
	if p.Total == 0 {
		return 0
	}
	return p.Counts[len(p.Counts)-1] / p.Total
}

func (p Profile) RateGap() float64 {
	if len(p.Rates) < 2 {
		return 0
	}
	maxGap := 0.0
	for i := 1; i < len(p.Rates); i++ {
		g := p.Rates[i] - p.Rates[i-1]
		if g < 0 {
			g = -g
		}
		if g > maxGap {
			maxGap = g
		}
	}
	return maxGap
}
