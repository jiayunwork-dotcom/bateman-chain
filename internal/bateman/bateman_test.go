package bateman

import (
	"errors"
	"math"
	"testing"

	"bateman-chain/internal/chain"
)

func TestTimeZeroReturnsInitialCounts(t *testing.T) {
	lambda := []float64{0.3, 1.2, 0}
	n0 := []float64{7, 3, 0}
	ev, err := New(lambda, n0)
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	counts, decayed, err := ev.At(0)
	if err != nil {
		t.Fatalf("At(0): expected no error, got %v", err)
	}
	for i := range n0 {
		if counts[i] != n0[i] {
			t.Errorf("nucleus %d at t=0: expected %g, got %g", i, n0[i], counts[i])
		}
	}
	if decayed != 0 {
		t.Errorf("decayed at t=0: expected 0, got %g", decayed)
	}
	act, err := ev.ActivityAt(0)
	if err != nil {
		t.Fatalf("ActivityAt(0): expected no error, got %v", err)
	}
	if math.Abs(act[0]-2.1) > 1e-12 || math.Abs(act[1]-3.6) > 1e-12 || act[2] != 0 {
		t.Errorf("activity at t=0: expected [2.1 3.6 0], got %v", act)
	}
}

func TestBatemanDenominatorUsesDifference(t *testing.T) {
	lambda := []float64{1.0, 0.2}
	n0 := []float64{100, 0}
	ev, err := New(lambda, n0)
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	counts, _, err := ev.At(2)
	if err != nil {
		t.Fatalf("At(2): expected no error, got %v", err)
	}
	wantN2 := 100 * (math.Exp(-2) - math.Exp(-0.4)) / (0.2 - 1)
	if math.Abs(counts[1]-wantN2) > 1e-6 {
		t.Errorf("daughter at t=2: expected %g, got %g", wantN2, counts[1])
	}
	if counts[1] <= 0 {
		t.Errorf("daughter must stay non-negative, got %g", counts[1])
	}
	wantN1 := 100 * math.Exp(-2)
	if math.Abs(counts[0]-wantN1) > 1e-9 {
		t.Errorf("parent at t=2: expected %g, got %g", wantN1, counts[0])
	}
	peak := FirstDaughterPeakTime(1.0, 0.2)
	countsPeak, _, err := ev.At(peak)
	if err != nil {
		t.Fatalf("At(peak): expected no error, got %v", err)
	}
	wantPeak := FirstDaughterPeakValue(100, 1.0, 0.2)
	if math.Abs(countsPeak[1]-wantPeak) > 1e-4 {
		t.Errorf("daughter at peak time %g: expected %g, got %g", peak, wantPeak, countsPeak[1])
	}
}

func TestParentCountsFollowExponential(t *testing.T) {
	lambda := []float64{0.3, 1.7, 0.9}
	n0 := []float64{42, 0, 0}
	ev, err := New(lambda, n0)
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	for _, tt := range []float64{0.5, 2, 5} {
		counts, _, err := ev.At(tt)
		if err != nil {
			t.Fatalf("At(%g): expected no error, got %v", tt, err)
		}
		want := 42 * math.Exp(-0.3*tt)
		if math.Abs(counts[0]-want) > 1e-9*math.Max(1, want) {
			t.Errorf("parent at t=%g: expected %g, got %g", tt, want, counts[0])
		}
	}
}

func TestDaughterLinearGrowthSmallTime(t *testing.T) {
	lambda := []float64{1e-4, 1e-3}
	n0 := []float64{1e6, 0}
	ev, err := New(lambda, n0)
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	counts, _, err := ev.At(10)
	if err != nil {
		t.Fatalf("At(10): expected no error, got %v", err)
	}
	linear := lambda[0] * n0[0] * 10
	if math.Abs(counts[1]-linear) > 0.01*linear {
		t.Errorf("daughter at small t: expected ~%g (lambda1*N1*t), got %g", linear, counts[1])
	}
}

func TestStableLastNuclideConservesTotal(t *testing.T) {
	lambda := []float64{1.0, 0.5, 0}
	n0 := []float64{100, 0, 0}
	ev, err := New(lambda, n0)
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	counts, decayed, err := ev.At(50)
	if err != nil {
		t.Fatalf("At(50): expected no error, got %v", err)
	}
	total := chain.Total(counts)
	if math.Abs(total-100) > 1e-6 {
		t.Errorf("stable terminal total: expected 100, got %g", total)
	}
	if counts[2] < 99.9 {
		t.Errorf("stable terminal: expected atoms accumulated on N3, got N3=%g", counts[2])
	}
	if decayed != 0 {
		t.Errorf("stable terminal decayed: expected 0, got %g", decayed)
	}
}

func TestAllUnstableVanishesAtLongTime(t *testing.T) {
	lambda := []float64{1.0, 0.5, 0.2}
	n0 := []float64{100, 0, 0}
	ev, err := New(lambda, n0)
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	counts, decayed, err := ev.At(200)
	if err != nil {
		t.Fatalf("At(200): expected no error, got %v", err)
	}
	total := chain.Total(counts)
	if total > 1e-4 {
		t.Errorf("all-unstable long time: expected total near 0, got %g", total)
	}
	if math.Abs(total+decayed-100) > 1e-6 {
		t.Errorf("decayed accounting: expected total+decayed=100, got %g", total+decayed)
	}
}

func TestShortLivedDaughterReachesQuasiSteady(t *testing.T) {
	lambda := []float64{1e-3, 1.0}
	n0 := []float64{1000, 0}
	ev, err := New(lambda, n0)
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	tt := 50.0
	counts, _, err := ev.At(tt)
	if err != nil {
		t.Fatalf("At(%g): expected no error, got %v", tt, err)
	}
	parent := 1000 * math.Exp(-1e-3*tt)
	want := 1e-3 * parent / 1.0
	if math.Abs(counts[1]-want) > 0.01*want {
		t.Errorf("short-lived daughter: expected ~%g, got %g", want, counts[1])
	}
}

func TestDecayedCountConservesAtoms(t *testing.T) {
	lambda := []float64{0.3, 0.7, 0.2}
	n0 := []float64{60, 25, 5}
	ev, err := New(lambda, n0)
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	for _, tt := range []float64{1, 5, 20} {
		counts, decayed, err := ev.At(tt)
		if err != nil {
			t.Fatalf("At(%g): expected no error, got %v", tt, err)
		}
		total := chain.Total(counts)
		if math.Abs(total+decayed-90) > 1e-6 {
			t.Errorf("at t=%g: expected total+decayed=90, got %g", tt, total+decayed)
		}
	}
}

func TestActivityIntegralMatchesReference(t *testing.T) {
	lambda := []float64{0.01, 0.02}
	n0 := []float64{1000, 0}
	ev, err := New(lambda, n0)
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	tt := 50.0
	want := TwoNucleusIntegral(tt, 1000, 0.01, 0.02)
	got, err := ev.Integral(1, tt)
	if err != nil {
		t.Fatalf("Integral: expected no error, got %v", err)
	}
	if math.Abs(got-want) > 1e-6*math.Max(1, want) {
		t.Errorf("daughter integral: expected %g, got %g", want, got)
	}
}

func TestNearEqualLambdaDetected(t *testing.T) {
	if !NearEqual([]float64{1.0, 1.0 + 1e-9}) {
		t.Error("near-equal pair: expected NearEqual true")
	}
	if !NearEqual([]float64{0, 0}) {
		t.Error("two stable nuclides: expected NearEqual true")
	}
	if NearEqual([]float64{1.0, 2.0, 0.5}) {
		t.Error("separated rates: expected NearEqual false")
	}
}

func TestEqualRateReferenceMatchesLimit(t *testing.T) {
	for _, tt := range []float64{1, 5, 20} {
		got := EqualRateSeries(tt, 100, 0.5, 2)
		want := 100 * math.Pow(0.5*tt, 2) / 2 * math.Exp(-0.5*tt)
		if math.Abs(got-want) > 1e-12*math.Max(1, want) {
			t.Errorf("equal-rate series at t=%g: expected %g, got %g", tt, want, got)
		}
	}
}

func TestVerifySeriesFlagsBrokenConservation(t *testing.T) {
	lambda := []float64{0.3, 0.7}
	n0 := []float64{50, 0}
	ev, err := New(lambda, n0)
	if err != nil {
		t.Fatalf("New: expected no error, got %v", err)
	}
	ser, err := ev.Series([]float64{0, 1, 5})
	if err != nil {
		t.Fatalf("Series: expected no error, got %v", err)
	}
	if err := VerifySeries(n0, ser, 1e-9); err != nil {
		t.Errorf("VerifySeries on a valid series: expected nil, got %v", err)
	}
	broken := ser
	broken.Rows[1] = []float64{1000, 0}
	if err := VerifySeries(n0, broken, 1e-9); !errors.Is(err, chain.ErrConservation) {
		t.Errorf("VerifySeries on a broken series: expected ErrConservation, got %v", err)
	}
}
