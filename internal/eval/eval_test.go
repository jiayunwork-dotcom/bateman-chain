package eval

import (
	"errors"
	"math"
	"strings"
	"testing"

	"bateman-chain/internal/chain"
)

func TestWellSeparatedCaseUsesBateman(t *testing.T) {
	c := Case{Lambda: []float64{1.0, 0.2}, Initial: []float64{100, 0}}
	req, err := Single(2)
	if err != nil {
		t.Fatalf("Single: expected no error, got %v", err)
	}
	res, err := Run(c, req)
	if err != nil {
		t.Fatalf("Run: expected no error, got %v", err)
	}
	if res.Solver != SolverBateman {
		t.Errorf("solver: expected %q, got %q", SolverBateman, res.Solver)
	}
	if math.Abs(res.Counts[0][1]-66.8731) > 1e-3 {
		t.Errorf("daughter at t=2: expected 66.8731, got %g", res.Counts[0][1])
	}
}

func TestNearEqualCaseRoutesToRK4(t *testing.T) {
	c := Case{Lambda: []float64{1.0, 1.0 + 1e-9}, Initial: []float64{10, 0}}
	req, err := Single(2)
	if err != nil {
		t.Fatalf("Single: expected no error, got %v", err)
	}
	res, err := Run(c, req)
	if err != nil {
		t.Fatalf("Run: expected no error, got %v", err)
	}
	if res.Solver != SolverRK4 {
		t.Errorf("solver: expected %q, got %q", SolverRK4, res.Solver)
	}
	if !chain.AllFinite(res.Counts[0]) {
		t.Errorf("counts: expected finite, got %v", res.Counts[0])
	}
	want := 10 * 2 * math.Exp(-2)
	if math.Abs(res.Counts[0][1]-want) > 1e-3 {
		t.Errorf("daughter limit at t=2: expected %g, got %g", want, res.Counts[0][1])
	}
}

func TestGridStepLimitErrors(t *testing.T) {
	if _, err := Grid(0, 10, 100001, 0); !errors.Is(err, chain.ErrStepLimit) {
		t.Errorf("oversized grid: expected ErrStepLimit, got %v", err)
	}
	if _, err := Grid(0, 10, 100, 50); !errors.Is(err, chain.ErrStepLimit) {
		t.Errorf("custom step limit: expected ErrStepLimit, got %v", err)
	}
}

func TestRejectsInvalidCaseAtEvalLevel(t *testing.T) {
	c := Case{Lambda: []float64{-0.5, 0.2}, Initial: []float64{100, 0}}
	req, err := Single(1)
	if err != nil {
		t.Fatalf("Single: expected no error, got %v", err)
	}
	if _, err := Run(c, req); !errors.Is(err, chain.ErrNegativeLambda) {
		t.Errorf("negative lambda: expected ErrNegativeLambda, got %v", err)
	}
}

func TestReportFieldsMatchSolver(t *testing.T) {
	c := Case{Names: []string{"A", "B"}, Lambda: []float64{0.3, 0.7}, Initial: []float64{40, 10}}
	req, err := Grid(0, 10, 3, 0)
	if err != nil {
		t.Fatalf("Grid: expected no error, got %v", err)
	}
	res, err := Run(c, req)
	if err != nil {
		t.Fatalf("Run: expected no error, got %v", err)
	}
	if len(res.Times) != 3 || len(res.Counts) != 3 {
		t.Fatalf("rows: expected 3, got %d", len(res.Counts))
	}
	for r := range res.Times {
		for i := 0; i < 2; i++ {
			if math.Abs(res.Activities[r][i]-res.Spec.Lambda[i]*res.Counts[r][i]) > 1e-12 {
				t.Errorf("row %d nuclide %d: A must equal lambda*N", r, i)
			}
		}
		if math.Abs(res.Totals[r]-res.Counts[r][0]-res.Counts[r][1]) > 1e-12 {
			t.Errorf("row %d: total must equal sum of counts", r)
		}
	}
	table := res.Table()
	if !strings.Contains(table, "A") || !strings.Contains(table, "solver") {
		t.Errorf("table: expected nuclide names and solver line, got:\n%s", table)
	}
}

func TestCaseJSONParseRejectsUnknown(t *testing.T) {
	if _, err := LoadCase([]byte(`{"lambda":[1,0.2],"initial":[100,0],"names":["a","b"],"nuclex":1}`)); err == nil {
		t.Error("unknown field: expected a parse error, got nil")
	}
	c, err := LoadCase([]byte(`{"lambda":[1,0.2],"initial":[100,0]}`))
	if err != nil {
		t.Fatalf("LoadCase: expected no error, got %v", err)
	}
	if len(c.Lambda) != 2 || c.Initial[0] != 100 {
		t.Errorf("parsed case: expected lambda [1 0.2] initial [100 0], got %v %v", c.Lambda, c.Initial)
	}
}

func TestCrossCheckClosedNumericAgree(t *testing.T) {
	c := Case{Lambda: []float64{0.3, 0.7, 0.2}, Initial: []float64{50, 5, 0}}
	cc, err := CrossCheckClosedNumeric(c, []float64{0.5, 2, 8}, 1e-3, 0)
	if err != nil {
		t.Fatalf("CrossCheckClosedNumeric: expected no error, got %v", err)
	}
	if !cc.WithinTol {
		t.Errorf("cross check: expected within tolerance, got %s", cc.Describe())
	}
}

func TestSummarizeReportsFinalState(t *testing.T) {
	c := Case{Lambda: []float64{0.3, 0.7}, Initial: []float64{40, 10}}
	req, err := Single(5)
	if err != nil {
		t.Fatalf("Single: expected no error, got %v", err)
	}
	res, err := Run(c, req)
	if err != nil {
		t.Fatalf("Run: expected no error, got %v", err)
	}
	sum := Summarize(res)
	if math.Abs(sum.ConservationErr) > 1e-6 {
		t.Errorf("conservation residual: expected ~0, got %g", sum.ConservationErr)
	}
	if math.Abs(sum.FinalTotal-(40*math.Exp(-1.5)+10*math.Exp(-3.5)+batemanDaughter(5))) > 1e-6 {
		t.Errorf("final total does not match the closed form")
	}
}

func batemanDaughter(t float64) float64 {
	return 40 * 0.3 * (math.Exp(-0.3*t) - math.Exp(-0.7*t)) / (0.7 - 0.3)
}

func TestActivityTotalsNonNegative(t *testing.T) {
	c := Case{Lambda: []float64{0.3, 0.7}, Initial: []float64{40, 10}}
	req, err := Grid(0, 10, 5, 0)
	if err != nil {
		t.Fatalf("Grid: expected no error, got %v", err)
	}
	res, err := Run(c, req)
	if err != nil {
		t.Fatalf("Run: expected no error, got %v", err)
	}
	totals := ActivityTotals(res)
	for r, tt := range totals {
		if tt < 0 {
			t.Errorf("row %d: activity total expected non-negative, got %g", r, tt)
		}
	}
}
