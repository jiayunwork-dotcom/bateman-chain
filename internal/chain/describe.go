package chain

import (
	"fmt"
	"strings"
)

func (s Spec) Header() string {
	names := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		names[i] = s.Name(i)
	}
	return strings.Join(names, " -> ")
}

func (s Spec) LastStable() bool {
	if s.Len() == 0 {
		return false
	}
	return s.Lambda[s.Len()-1] == 0
}

func (s Spec) AllUnstable() bool {
	for i := 0; i < s.Len(); i++ {
		if s.Lambda[i] == 0 {
			return false
		}
	}
	return true
}

func (s Spec) MaxLambda() float64 {
	max := 0.0
	for i := 0; i < s.Len(); i++ {
		if s.Lambda[i] > max {
			max = s.Lambda[i]
		}
	}
	return max
}

func (s Spec) Summary() string {
	parts := []string{
		fmt.Sprintf("length=%d", s.Len()),
		fmt.Sprintf("initial_total=%g", InitialTotal(s)),
		fmt.Sprintf("parent_activity=%g", s.Lambda[0]*s.N0[0]),
	}
	return strings.Join(parts, ", ")
}
