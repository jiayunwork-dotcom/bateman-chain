package eval

import (
	"fmt"
	"strings"

	"bateman-chain/internal/chain"
)

func (r Result) Table() string {
	var b strings.Builder
	fmt.Fprintf(&b, "chain: %s\n", r.Spec.Describe())
	fmt.Fprintf(&b, "half-lives: %s\n", r.halfLifeLine())
	fmt.Fprintf(&b, "solver: %s\n", r.Solver)
	for row, t := range r.Times {
		fmt.Fprintf(&b, "\nt = %g s\n", t)
		fmt.Fprintf(&b, "%-14s %-18s %-18s\n", "nuclide", "N", "A")
		for i := 0; i < r.Spec.Len(); i++ {
			fmt.Fprintf(&b, "%-14s %-18.6e %-18.6e\n",
				r.Spec.Name(i), r.Counts[row][i], r.Activities[row][i])
		}
		fmt.Fprintf(&b, "%-14s %-18.6e\n", "total", r.Totals[row])
		fmt.Fprintf(&b, "%-14s %-18.6e\n", "decayed", r.Decayed[row])
	}
	return b.String()
}

func (r Result) halfLifeLine() string {
	parts := make([]string, r.Spec.Len())
	for i := 0; i < r.Spec.Len(); i++ {
		parts[i] = fmt.Sprintf("%s=%.3e", r.Spec.Name(i), chain.HalfLife(r.Spec.Lambda[i]))
	}
	return strings.Join(parts, ", ")
}
