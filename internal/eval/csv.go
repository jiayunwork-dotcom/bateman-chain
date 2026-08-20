package eval

import (
	"fmt"
	"strings"
)

func (r Result) CSV() string {
	var b strings.Builder
	b.WriteString("time,index,nuclide,n,a,total,decayed\n")
	for row, t := range r.Times {
		for i := 0; i < r.Spec.Len(); i++ {
			fmt.Fprintf(&b, "%g,%d,%s,%g,%g,%g,%g\n",
				t, i, r.Spec.Name(i), r.Counts[row][i], r.Activities[row][i],
				r.Totals[row], r.Decayed[row])
		}
	}
	return b.String()
}
