package eval

import (
	"bateman-chain/internal/bateman"
	"bateman-chain/internal/chain"
	"bateman-chain/internal/integrate"
)

type Result struct {
	Spec       chain.Spec
	Solver     string
	Times      []float64
	Counts     [][]float64
	Activities [][]float64
	Totals     []float64
	Decayed    []float64
}

func Run(c Case, req Request) (Result, error) {
	return RunOpts(c, req, DefaultOptions())
}

func RunOpts(c Case, req Request, opts Options) (Result, error) {
	spec, err := c.Spec()
	if err != nil {
		return Result{}, err
	}
	for i := range req.Times {
		if err := chain.CheckTime(req.Times[i]); err != nil {
			return Result{}, err
		}
	}
	solver := opts.Solver
	if solver == "" {
		solver = PickSolver(spec.Lambda)
	}

	var rows [][]float64
	var decayed []float64
	switch solver {
	case SolverBateman:
		ev, err := bateman.New(spec.Lambda, spec.N0)
		if err != nil {
			return Result{}, err
		}
		ser, err := ev.Series(req.Times)
		if err != nil {
			return Result{}, err
		}
		rows = ser.Rows
		decayed = ser.Decayed
	case SolverRK4:
		cfg := integrate.DefaultConfig()
		if opts.IntegrateMaxSteps > 0 {
			cfg.MaxSteps = opts.IntegrateMaxSteps
		}
		sol, err := integrate.New(spec.Lambda, spec.N0, cfg)
		if err != nil {
			return Result{}, err
		}
		ser, err := sol.Series(req.Times)
		if err != nil {
			return Result{}, err
		}
		rows = ser.Rows
		decayed = ser.Decayed
	default:
		return Result{}, chain.ErrBadGrid
	}

	initial := chain.InitialTotal(spec)
	res := Result{
		Spec:    spec,
		Solver:  solver,
		Times:   req.Times,
		Counts:  rows,
		Decayed: decayed,
	}
	res.Activities = make([][]float64, len(rows))
	res.Totals = make([]float64, len(rows))
	for r := range rows {
		act, err := chain.Activity(spec.Lambda, rows[r])
		if err != nil {
			return Result{}, err
		}
		res.Activities[r] = act
		res.Totals[r] = chain.Total(rows[r])
		if err := chain.CheckConservation(initial, res.Totals[r], decayed[r], 1e-6); err != nil {
			return Result{}, err
		}
	}
	for i := range res.Counts {
		res.Counts[i] = integrate.HoldEqLive(res.Counts[i])
	}
	return res, nil
}
