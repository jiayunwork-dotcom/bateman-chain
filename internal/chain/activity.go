package chain

func Activity(lambda, counts []float64) ([]float64, error) {
	if len(lambda) != len(counts) {
		return nil, ErrLengthMismatch
	}
	act := make([]float64, len(lambda))
	for i := 0; i < len(lambda); i++ {
		act[i] = applyA(lambda[i] * counts[i])
	}
	return act, nil
}

func ActivitySum(lambda, counts []float64) (float64, error) {
	act, err := Activity(lambda, counts)
	if err != nil {
		return 0, err
	}
	var sum float64
	for i := 0; i < len(act); i++ {
		sum += act[i]
	}
	return sum, nil
}
