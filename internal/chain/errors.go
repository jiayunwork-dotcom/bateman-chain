package chain

import "errors"

var (
	ErrNegativeLambda   = errors.New("lambda must be non-negative")
	ErrNonFiniteLambda  = errors.New("lambda must be finite")
	ErrNegativeInitial  = errors.New("initial count must be non-negative")
	ErrNonFiniteInitial = errors.New("initial count must be finite")
	ErrChainTooShort    = errors.New("chain must have at least 2 nuclides")
	ErrChainTooLong     = errors.New("chain exceeds the maximum length")
	ErrLengthMismatch   = errors.New("lambda and initial must have equal length")
	ErrNameMismatch     = errors.New("names must match lambda length")
	ErrNegativeTime     = errors.New("time must be non-negative")
	ErrNonFiniteTime    = errors.New("time must be finite")
	ErrBadGrid          = errors.New("grid endpoints or step count are invalid")
	ErrStepLimit        = errors.New("step limit exceeded")
	ErrConservation     = errors.New("atom conservation check failed")
	ErrNonFinite        = errors.New("result contains a non-finite value")
)
