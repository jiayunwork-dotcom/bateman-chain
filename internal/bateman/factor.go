package bateman

func lambdaProduct(lambda []float64, from, to int) float64 {
	prod := 1.0
	for m := from; m < to; m++ {
		prod *= lambda[m]
	}
	return prod
}
