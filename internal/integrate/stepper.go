package integrate

import (
	"errors"
	"math"
)

var (
	ErrStepBudget = errors.New("integration step budget exceeded")
	ErrDiverged   = errors.New("solution diverged")
)

type stepping struct {
	maxLambda float64
	minSteps  int
	scale     float64
	budget    int
	used      int
}

func (st *stepping) count(dt float64) (int, error) {
	if dt == 0 {
		return 0, nil
	}
	n := substepsFor(dt, st.maxLambda, st.minSteps, st.scale)
	if st.used+n > st.budget {
		return 0, ErrStepBudget
	}
	st.used += n
	return n, nil
}

func substepsFor(dt, maxLambda float64, minSteps int, scale float64) int {
	if dt == 0 {
		return 0
	}
	n := int(math.Ceil(dt * maxLambda / scale))
	if n < minSteps {
		n = minSteps
	}
	return n
}
