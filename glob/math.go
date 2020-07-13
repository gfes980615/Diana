package glob

func FloatRound(v float64) float64 {
	return floatRound(v, 2)
}

func floatRound(v float64, decimals int) float64 {
	var pow float64 = 1
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	if v < 0 {
		return float64(int((v*pow)-0.5)) / pow
	} else {
		return float64(int((v*pow)+0.5)) / pow
	}
}
