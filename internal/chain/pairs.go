package chain

import "sort"

type RateEntry struct {
	Index int
	Rate  float64
	Name  string
}

func (s Spec) RateOrder() []RateEntry {
	entries := make([]RateEntry, s.Len())
	for i := 0; i < s.Len(); i++ {
		entries[i] = RateEntry{Index: i, Rate: s.Lambda[i], Name: s.Name(i)}
	}
	sort.SliceStable(entries, func(a, b int) bool {
		return entries[a].Rate < entries[b].Rate
	})
	return entries
}

func (s Spec) SlowestNuclide() int {
	idx := 0
	for i := 1; i < s.Len(); i++ {
		if s.Lambda[i] < s.Lambda[idx] {
			idx = i
		}
	}
	return idx
}

func (s Spec) FastestNuclide() int {
	idx := 0
	for i := 1; i < s.Len(); i++ {
		if s.Lambda[i] > s.Lambda[idx] {
			idx = i
		}
	}
	return idx
}
