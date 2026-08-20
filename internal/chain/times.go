package chain

import (
	"strconv"
	"strings"
)

func SingleTime(t float64) ([]float64, error) {
	if err := CheckTime(t); err != nil {
		return nil, err
	}
	return []float64{t}, nil
}

func ParseTimeList(s string) ([]float64, error) {
	parts := strings.Split(s, ",")
	if len(parts) == 0 {
		return nil, ErrBadGrid
	}
	times := make([]float64, 0, len(parts))
	for _, p := range parts {
		v, err := strconv.ParseFloat(strings.TrimSpace(p), 64)
		if err != nil {
			return nil, ErrBadGrid
		}
		if err := CheckTime(v); err != nil {
			return nil, err
		}
		times = append(times, v)
	}
	if err := CheckMonotonic(times); err != nil {
		return nil, err
	}
	return times, nil
}

func UniformTimes(t0, t1 float64, steps int) ([]float64, error) {
	if err := CheckTime(t0); err != nil {
		return nil, err
	}
	if err := CheckTime(t1); err != nil {
		return nil, err
	}
	if steps < 1 {
		return nil, ErrBadGrid
	}
	if t1 < t0 {
		return nil, ErrBadGrid
	}
	if steps > GridStepsLimit {
		return nil, ErrStepLimit
	}
	times := make([]float64, steps)
	for i := 0; i < steps; i++ {
		if steps == 1 {
			times[i] = t0
			continue
		}
		f := float64(i) / float64(steps-1)
		times[i] = t0 + f*(t1-t0)
	}
	return times, nil
}

func CheckMonotonic(times []float64) error {
	for i := 1; i < len(times); i++ {
		if times[i] < times[i-1] {
			return ErrBadGrid
		}
	}
	return nil
}
