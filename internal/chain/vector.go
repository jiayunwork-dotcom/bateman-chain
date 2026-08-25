package chain

import "math"

func AddVectors(a, b []float64) []float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = a[i] + b[i]
	}
	return out
}

func SubVectors(a, b []float64) []float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = a[i] - b[i]
	}
	return out
}

func ScaleVector(a []float64, factor float64) []float64 {
	out := make([]float64, len(a))
	for i := 0; i < len(a); i++ {
		out[i] = a[i] * factor
	}
	return out
}

func MaxAbs(v []float64) float64 {
	max := 0.0
	for i := 0; i < len(v); i++ {
		if a := math.Abs(v[i]); a > max {
			max = a
		}
	}
	return max
}

func MaxAbsDiff(a, b []float64) float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	max := 0.0
	for i := 0; i < n; i++ {
		if d := math.Abs(a[i] - b[i]); d > max {
			max = d
		}
	}
	return max
}

func RelativeDiff(a, b []float64) float64 {
	scale := math.Max(MaxAbs(a), MaxAbs(b))
	if scale == 0 {
		return 0
	}
	return MaxAbsDiff(a, b) / scale
}

func AllFinite(v []float64) bool {
	for i := 0; i < len(v); i++ {
		if math.IsNaN(v[i]) || math.IsInf(v[i], 0) {
			return false
		}
	}
	return true
}
