package chain

import (
	"errors"
	"math"
	"testing"
)

func TestRejectsInvalidChainInputs(t *testing.T) {
	cases := []struct {
		name    string
		lambda  []float64
		n0      []float64
		wantErr error
	}{
		{name: "negative lambda", lambda: []float64{-0.5, 0.2}, n0: []float64{100, 0}, wantErr: ErrNegativeLambda},
		{name: "single nuclide", lambda: []float64{0.5}, n0: []float64{100}, wantErr: ErrChainTooShort},
	}
	for _, tc := range cases {
		spec, err := NewSpec(tc.lambda, tc.n0, nil)
		if !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: expected error %v, got %v", tc.name, tc.wantErr, err)
		}
		if spec.Len() != 0 && err == nil {
			t.Errorf("%s: expected the invalid spec to be rejected", tc.name)
		}
	}
	if err := CheckTime(-1); !errors.Is(err, ErrNegativeTime) {
		t.Errorf("negative time: expected ErrNegativeTime, got %v", err)
	}
	if err := CheckTime(math.Inf(1)); !errors.Is(err, ErrNonFiniteTime) {
		t.Errorf("infinite time: expected ErrNonFiniteTime, got %v", err)
	}
}

func TestRejectsChainTooLong(t *testing.T) {
	lambda := make([]float64, MaxChainLen+1)
	n0 := make([]float64, MaxChainLen+1)
	for i := range lambda {
		lambda[i] = 0.1
		n0[i] = 1
	}
	if _, err := NewSpec(lambda, n0, nil); !errors.Is(err, ErrChainTooLong) {
		t.Errorf("13-nuclide chain: expected ErrChainTooLong, got %v", err)
	}
}

func TestRejectsNegativeInitialCounts(t *testing.T) {
	if _, err := NewSpec([]float64{0.5, 0.2}, []float64{100, -1}, nil); !errors.Is(err, ErrNegativeInitial) {
		t.Errorf("negative initial count: expected ErrNegativeInitial, got %v", err)
	}
}

func TestRejectsMismatchedLengths(t *testing.T) {
	if _, err := NewSpec([]float64{0.5, 0.2, 0.1}, []float64{100, 0}, nil); !errors.Is(err, ErrLengthMismatch) {
		t.Errorf("length mismatch: expected ErrLengthMismatch, got %v", err)
	}
	if _, err := NewSpec([]float64{0.5, 0.2}, []float64{100, 0}, []string{"a"}); !errors.Is(err, ErrNameMismatch) {
		t.Errorf("name mismatch: expected ErrNameMismatch, got %v", err)
	}
}

func TestUniformGridBuildsMonotone(t *testing.T) {
	times, err := UniformTimes(0, 10, 5)
	if err != nil {
		t.Fatalf("UniformTimes: expected no error, got %v", err)
	}
	want := []float64{0, 2.5, 5, 7.5, 10}
	for i := range want {
		if times[i] != want[i] {
			t.Errorf("grid[%d]: expected %g, got %g", i, want[i], times[i])
		}
	}
	if _, err := UniformTimes(5, 0, 3); !errors.Is(err, ErrBadGrid) {
		t.Errorf("reversed grid: expected ErrBadGrid, got %v", err)
	}
	if _, err := UniformTimes(0, 10, GridStepsLimit+1); !errors.Is(err, ErrStepLimit) {
		t.Errorf("oversized grid: expected ErrStepLimit, got %v", err)
	}
}

func TestParseTimeListAcceptsCommaSeparated(t *testing.T) {
	times, err := ParseTimeList("0, 10, 100")
	if err != nil {
		t.Fatalf("ParseTimeList: expected no error, got %v", err)
	}
	if len(times) != 3 || times[0] != 0 || times[1] != 10 || times[2] != 100 {
		t.Errorf("parsed times: expected [0 10 100], got %v", times)
	}
	if _, err := ParseTimeList("10, 0"); !errors.Is(err, ErrBadGrid) {
		t.Errorf("decreasing list: expected ErrBadGrid, got %v", err)
	}
	if _, err := ParseTimeList("abc"); !errors.Is(err, ErrBadGrid) {
		t.Errorf("garbage list: expected ErrBadGrid, got %v", err)
	}
}

func TestActivityUsesLambdaVector(t *testing.T) {
	lambda := []float64{0.3, 1.2, 0}
	n0 := []float64{10, 5, 0}
	act, err := Activity(lambda, n0)
	if err != nil {
		t.Fatalf("Activity: expected no error, got %v", err)
	}
	if act[0] != 3.0 || act[1] != 6.0 || act[2] != 0 {
		t.Errorf("activity at initial counts: expected [3 6 0], got %v", act)
	}
	if _, err := Activity([]float64{1}, n0); !errors.Is(err, ErrLengthMismatch) {
		t.Errorf("activity length mismatch: expected ErrLengthMismatch, got %v", err)
	}
}

func TestConservationCheckToleratesRoundoff(t *testing.T) {
	if err := CheckConservation(100, 40, 60, 1e-12); err != nil {
		t.Errorf("exact conservation: expected nil, got %v", err)
	}
	if err := CheckConservation(100, 40, 59, 1e-12); !errors.Is(err, ErrConservation) {
		t.Errorf("broken conservation: expected ErrConservation, got %v", err)
	}
}

func TestHalfLifeAndLifetime(t *testing.T) {
	if HalfLife(0) != math.Inf(1) {
		t.Errorf("stable nuclide half-life: expected +Inf, got %g", HalfLife(0))
	}
	if got := HalfLife(math.Ln2); math.Abs(got-1) > 1e-12 {
		t.Errorf("half-life of Ln2: expected 1, got %g", got)
	}
	if MeanLifetime(2) != 0.5 {
		t.Errorf("mean lifetime of rate 2: expected 0.5, got %g", MeanLifetime(2))
	}
}

func TestRateOrderAndExtremes(t *testing.T) {
	spec, err := NewSpec([]float64{0.5, 0.1, 2.0}, []float64{1, 1, 1}, nil)
	if err != nil {
		t.Fatalf("NewSpec: expected no error, got %v", err)
	}
	order := spec.RateOrder()
	if order[0].Index != 1 || order[1].Index != 0 || order[2].Index != 2 {
		t.Errorf("rate order: expected indices [1 0 2], got [%d %d %d]", order[0].Index, order[1].Index, order[2].Index)
	}
	if spec.SlowestNuclide() != 1 {
		t.Errorf("slowest nuclide: expected index 1, got %d", spec.SlowestNuclide())
	}
	if spec.FastestNuclide() != 2 {
		t.Errorf("fastest nuclide: expected index 2, got %d", spec.FastestNuclide())
	}
}
