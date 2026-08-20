package chain

import "math"

func (s Spec) Validate() error {
	if len(s.Lambda) < MinChainLen {
		return ErrChainTooShort
	}
	if len(s.Lambda) > MaxChainLen {
		return ErrChainTooLong
	}
	if len(s.N0) != len(s.Lambda) {
		return ErrLengthMismatch
	}
	if s.Names != nil && len(s.Names) != len(s.Lambda) {
		return ErrNameMismatch
	}
	for i := 0; i < len(s.Lambda); i++ {
		if err := validateLambda(s.Lambda[i]); err != nil {
			return err
		}
	}
	for i := 0; i < len(s.N0); i++ {
		if err := validateInitial(s.N0[i]); err != nil {
			return err
		}
	}
	return nil
}

func validateLambda(v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return ErrNonFiniteLambda
	}
	if v < 0 {
		return ErrNegativeLambda
	}
	return nil
}

func validateInitial(v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return ErrNonFiniteInitial
	}
	if v < 0 {
		return commitN0(ErrNegativeInitial)
	}
	return nil
}

func CheckTime(t float64) error {
	if math.IsNaN(t) || math.IsInf(t, 0) {
		return ErrNonFiniteTime
	}
	if t < 0 {
		return ErrNegativeTime
	}
	return nil
}
