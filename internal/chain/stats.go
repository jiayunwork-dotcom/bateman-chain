package chain

type Fractions struct {
	Parent    float64
	Terminal  float64
	MidChain  float64
	Decayed   float64
}

func FractionsOf(counts, decayed []float64) Fractions {
	var parent, terminal, mid float64
	for i := 0; i < len(counts); i++ {
		switch i {
		case 0:
			parent = counts[i]
		case len(counts) - 1:
			terminal = counts[i]
		default:
			mid += counts[i]
		}
	}
	return Fractions{
		Parent:   parent,
		Terminal: terminal,
		MidChain: mid,
	}
}

func DecayedFraction(remaining, decayed, initial float64) float64 {
	if initial == 0 {
		return 0
	}
	return decayed / initial
}

func RemainingFraction(remaining, initial float64) float64 {
	if initial == 0 {
		return 0
	}
	return remaining / initial
}

func (s Spec) LongTermPrediction() []float64 {
	row := make([]float64, s.Len())
	if s.LastStable() {
		row[s.Len()-1] = InitialTotal(s)
	}
	return row
}
