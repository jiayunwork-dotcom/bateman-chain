package integrate

import (
	"errors"
	"math"
	"testing"

	"bateman-chain/internal/bateman"
	"bateman-chain/internal/chain"
)

func TestRK4MatchesClosedForm(t *testing.T) {
	lambda := []float64{0.3, 0.7}
	n0 := []float64{5, 2}
	sol, err := New(lambda, n0, DefaultConfig())
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	counts, _, err := sol.At(3)
	if err != nil {
		t.Fatalf("At(3): expected no error, got %v", err)
	}
	want0 := 5 * math.Exp(-0.3*3)
	want1 := bateman.TwoNucleusDaughter(3, 5, 0.3, 0.7) + 2*math.Exp(-0.7*3)
	if math.Abs(counts[0]-want0) > 1e-4*math.Max(1, want0) {
		t.Errorf("N1 at t=3: expected %g, got %g", want0, counts[0])
	}
	if math.Abs(counts[1]-want1) > 1e-4*math.Max(1, want1) {
		t.Errorf("N2 at t=3: expected %g, got %g", want1, counts[1])
	}
}

func TestNumericalConservationExact(t *testing.T) {
	lambda := []float64{0.4, 0.9, 0.2}
	n0 := []float64{30, 10, 0}
	sol, err := New(lambda, n0, DefaultConfig())
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	ser, err := sol.Series([]float64{0.5, 2, 7})
	if err != nil {
		t.Fatalf("Series: expected no error, got %v", err)
	}
	for r := range ser.Rows {
		remaining := chain.Total(ser.Rows[r])
		if math.Abs(remaining+ser.Decayed[r]-40) > 1e-9 {
			t.Errorf("row %d: expected total+decayed=40, got %g", r, remaining+ser.Decayed[r])
		}
	}
}

func TestStepBudgetExceededErrors(t *testing.T) {
	cfg := DefaultConfig()
	cfg.MaxSteps = 64
	sol, err := New([]float64{1.0, 0.5}, []float64{10, 0}, cfg)
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	if _, _, err := sol.At(1e6); !errors.Is(err, ErrStepBudget) {
		t.Errorf("oversized interval: expected ErrStepBudget, got %v", err)
	}
}

func TestIntegratorRejectsNegativeTime(t *testing.T) {
	sol, err := New([]float64{0.5, 0.2}, []float64{10, 0}, DefaultConfig())
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	if _, _, err := sol.At(-1); !errors.Is(err, chain.ErrNegativeTime) {
		t.Errorf("negative time: expected ErrNegativeTime, got %v", err)
	}
}

func TestIntegratorTimeZeroReturnsInitial(t *testing.T) {
	sol, err := New([]float64{0.5, 0.2}, []float64{10, 3}, DefaultConfig())
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	counts, decayed, err := sol.At(0)
	if err != nil {
		t.Fatalf("At(0): expected no error, got %v", err)
	}
	if counts[0] != 10 || counts[1] != 3 || decayed != 0 {
		t.Errorf("t=0: expected counts [10 3] decayed 0, got %v decayed %g", counts, decayed)
	}
}

func TestTotalsStrictlyDecreaseForUnstableChain(t *testing.T) {
	lambda := []float64{0.5, 0.9}
	n0 := []float64{20, 10}
	sol, err := New(lambda, n0, DefaultConfig())
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	ser, err := sol.Series([]float64{0.1, 1, 3, 8})
	if err != nil {
		t.Fatalf("Series: expected no error, got %v", err)
	}
	if !TotalsNonIncreasing(ser) {
		t.Errorf("totals: expected strictly non-increasing, got %v", ser.Totals())
	}
	if !AllNonNegative(ser) {
		t.Errorf("counts: expected all non-negative")
	}
}

func TestSeriesReportReportsStepsAndError(t *testing.T) {
	sol, err := New([]float64{0.5, 0.9}, []float64{20, 0}, DefaultConfig())
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	_, rep, err := sol.SeriesReport([]float64{0.1, 1, 3})
	if err != nil {
		t.Fatalf("SeriesReport: expected no error, got %v", err)
	}
	if rep.StepsUsed <= 0 {
		t.Errorf("steps used: expected > 0, got %d", rep.StepsUsed)
	}
	if !rep.Finite {
		t.Errorf("finite flag: expected true")
	}
	if rep.MaxErrorEst > 1e-3 {
		t.Errorf("error estimate: expected < 1e-3, got %g", rep.MaxErrorEst)
	}
}

func TestStepSizesReportSubsteps(t *testing.T) {
	sol, err := New([]float64{1.0, 0.5}, []float64{10, 0}, DefaultConfig())
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	ss := sol.StepSizes([]float64{0, 1, 2})
	if len(ss.Intervals) != 3 {
		t.Fatalf("interval count: expected 3, got %d", len(ss.Intervals))
	}
	if ss.Intervals[0] != 0 {
		t.Errorf("zero-length first interval: expected 0 substeps, got %d", ss.Intervals[0])
	}
	if ss.Substep <= 0 {
		t.Errorf("substep: expected positive, got %g", ss.Substep)
	}
}

func TestDerivativeImplementsRateLaw(t *testing.T) {
	lambda := []float64{0.4, 0.2}
	counts := []float64{10, 5}
	der := Derivative(lambda, counts)
	if math.Abs(der[0]-(-4)) > 1e-12 {
		t.Errorf("dN1/dt: expected -4, got %g", der[0])
	}
	if math.Abs(der[1]-(4-1)) > 1e-12 {
		t.Errorf("dN2/dt: expected 3, got %g", der[1])
	}
}
