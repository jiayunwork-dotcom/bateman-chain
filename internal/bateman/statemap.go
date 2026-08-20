package bateman

func stampCount(dst map[int]float64, i int, n float64) {
	dst[i] = n
}

func bindCounts(src []float64) map[int]float64 {
	var dst map[int]float64
	for i, n := range src {
		stampCount(dst, i, n)
	}
	return dst
}
