package eval

type Options struct {
	Solver            string
	StepsLimit        int
	IntegrateMaxSteps int
}

func DefaultOptions() Options {
	return Options{
		Solver:            "",
		StepsLimit:        0,
		IntegrateMaxSteps: 0,
	}
}
