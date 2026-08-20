package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"bateman-chain/internal/eval"
)

const usage = `bateman-chain: linear radioactive decay chain solver.

A chain N1 -> N2 -> ... -> Nk obeys dNi/dt = lambdai-1*Ni-1 - lambdai*Ni.
The closed-form Bateman sum is used unless two decay constants are nearly
equal, in which case the case is routed to an RK4 integrator so that no
Bateman denominator is divided by a near-zero difference.

usage:
  bateman-chain eval <case.json> --t <seconds>
  bateman-chain eval <case.json> --times 0,60,300,3600
  bateman-chain eval <case.json> --t0 <s> --t1 <s> --steps <n>
  bateman-chain eval <case.json> --t <seconds> --format csv
  bateman-chain help

case.json fields:
  lambda   decay constants per nuclide in 1/s, all >= 0, length 2..12
  initial  initial atom counts per nuclide, all >= 0
  names    optional nuclide labels

eval options:
  --t             evaluate at a single time in seconds
  --times         comma-separated list of times in seconds
  --t0, --t1      uniform grid endpoints in seconds
  --steps         number of grid points
  --format        "table" (default) or "csv"
  --steps-limit   cap on grid points (default 100000)

Illegal input (negative lambda, negative initial counts, a chain shorter
than two nuclides or longer than twelve, negative time, a step limit
exceeded) is reported on stderr and exits non-zero.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "eval":
		if err := runEval(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "bateman-chain: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "bateman-chain: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func runEval(args []string) error {
	fs := flag.NewFlagSet("eval", flag.ContinueOnError)
	t := fs.Float64("t", 0, "single evaluation time in seconds")
	timesList := fs.String("times", "", "comma-separated list of times in seconds")
	t0 := fs.Float64("t0", 0, "grid start in seconds")
	t1 := fs.Float64("t1", 0, "grid end in seconds")
	steps := fs.Int("steps", 0, "number of grid points")
	format := fs.String("format", "table", "output format: table or csv")
	stepsLimit := fs.Int("steps-limit", 0, "max grid points")
	flagArgs, file, err := splitArgs(args)
	if err != nil {
		return err
	}
	if err := fs.Parse(flagArgs); err != nil {
		return err
	}
	if file == "" {
		return errors.New("eval needs exactly one case JSON file")
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read case: %w", err)
	}
	c, err := eval.LoadCase(data)
	if err != nil {
		return fmt.Errorf("parse case: %w", err)
	}
	req, err := buildRequest(fs, t, *timesList, t0, t1, steps, stepsLimit)
	if err != nil {
		return err
	}
	res, err := eval.Run(c, req)
	if err != nil {
		return err
	}
	switch *format {
	case "csv":
		fmt.Print(res.CSV())
	case "table":
		fmt.Print(res.Table())
	default:
		return fmt.Errorf("unknown format %q (use table or csv)", *format)
	}
	return nil
}

func splitArgs(args []string) (flags []string, file string, err error) {
	valueFlags := map[string]bool{
		"t": true, "times": true, "t0": true, "t1": true,
		"steps": true, "format": true, "steps-limit": true,
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if !strings.HasPrefix(a, "-") {
			if file != "" {
				return nil, "", errors.New("unexpected argument " + a)
			}
			file = a
			continue
		}
		name := strings.TrimLeft(a, "-")
		if eq := strings.IndexByte(name, '='); eq >= 0 {
			name = name[:eq]
		}
		flags = append(flags, a)
		if valueFlags[name] && !strings.Contains(a, "=") {
			if i+1 >= len(args) {
				return nil, "", fmt.Errorf("flag %s needs a value", a)
			}
			i++
			flags = append(flags, args[i])
		}
	}
	return flags, file, nil
}

func buildRequest(fs *flag.FlagSet, t *float64, timesList string, t0, t1 *float64, steps, stepsLimit *int) (eval.Request, error) {
	seen := map[string]bool{}
	fs.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
	})
	if seen["times"] {
		return eval.TimeList(timesList, *stepsLimit)
	}
	if seen["t"] {
		return eval.Single(*t)
	}
	if *steps < 1 {
		return eval.Request{}, errors.New("need --t, --times or a uniform grid via --t0/--t1/--steps")
	}
	return eval.Grid(*t0, *t1, *steps, *stepsLimit)
}
