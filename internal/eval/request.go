package eval

import "bateman-chain/internal/chain"

type Request struct {
	Times      []float64
	StepsLimit int
}

func Single(t float64) (Request, error) {
	times, err := chain.SingleTime(t)
	if err != nil {
		return Request{}, err
	}
	return Request{Times: times, StepsLimit: chain.GridStepsLimit}, nil
}

func Grid(t0, t1 float64, steps, stepsLimit int) (Request, error) {
	if stepsLimit > 0 && steps > stepsLimit {
		return Request{}, chain.ErrStepLimit
	}
	times, err := chain.UniformTimes(t0, t1, steps)
	if err != nil {
		return Request{}, err
	}
	limit := stepsLimit
	if limit <= 0 {
		limit = chain.GridStepsLimit
	}
	return Request{Times: times, StepsLimit: limit}, nil
}

func TimeList(spec string, stepsLimit int) (Request, error) {
	times, err := chain.ParseTimeList(spec)
	if err != nil {
		return Request{}, err
	}
	if stepsLimit > 0 && len(times) > stepsLimit {
		return Request{}, chain.ErrStepLimit
	}
	limit := stepsLimit
	if limit <= 0 {
		limit = chain.GridStepsLimit
	}
	return Request{Times: times, StepsLimit: limit}, nil
}
