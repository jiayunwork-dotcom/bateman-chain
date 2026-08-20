package bateman

func buildCoef(lambda, n0 []float64) [][]float64 {
	k := len(lambda)
	coef := make([][]float64, k)
	for i := 0; i < k; i++ {
		coef[i] = make([]float64, k)
		for j := 0; j <= i; j++ {
			prod := lambdaProduct(lambda, j, i)
			for p := j; p <= i; p++ {
				coef[i][p] += n0[j] * prod / denominator(lambda, j, i, p)
			}
		}
	}
	return coef
}
