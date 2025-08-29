package utils

// 计算 MACD：12EMA快线，26EMA慢线，9MACD信号，返回MACD集合，信号集合，柱子集合
func CalculateMACD(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) (macdLine, signalLine, histogram []float64) {
	emaFast := CalculateEMA(closePrices, fastPeriod)
	emaSlow := CalculateEMA(closePrices, slowPeriod)
	macdLine = make([]float64, len(closePrices))
	for i := range closePrices {
		macdLine[i] = emaFast[i] - emaSlow[i]
	}
	signalLine = CalculateEMA(macdLine, signalPeriod) //信号只是MACD的EMA平均
	histogram = make([]float64, len(closePrices))
	for i := range closePrices {
		histogram[i] = macdLine[i] - signalLine[i]
	}
	return
}

//为正
func IsGoldenUP(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	if len(closePrices) < slowPeriod+signalPeriod+1 {
		return false
	}

	_, _, histogram := CalculateMACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	if len(histogram) < 5 {
		return false
	}

	D := histogram[len(histogram)-2]
	E := histogram[len(histogram)-1]

	if E > 0 {
		return true
	}
	if D > 0 {
		return true
	}

	if D < 0 && E < 0 && D < E {
		return true
	}

	return false
}

//为正
func IsGolden(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	if len(closePrices) < slowPeriod+signalPeriod+1 {
		return false
	}

	_, _, histogram := CalculateMACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	if len(histogram) < 5 {
		return false
	}
	D := histogram[len(histogram)-2]
	E := histogram[len(histogram)-1]

	if E > 0 {
		return true
	}
	if D > 0 {
		return true
	}

	return false
}

//为负
func IsDeadDOWN(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	if len(closePrices) < slowPeriod+signalPeriod+1 {
		return false
	}

	_, _, histogram := CalculateMACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	if len(histogram) < 5 {
		return false
	}

	D := histogram[len(histogram)-2]
	E := histogram[len(histogram)-1]

	if E < 0 {
		return true
	}
	if D < 0 {
		return true
	}

	if D > 0 && E > 0 && D > E {
		return true
	}

	return false
}

// 判断是否为负
func IsDead(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	if len(closePrices) < slowPeriod+signalPeriod+1 {
		return false
	}

	_, _, histogram := CalculateMACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	if len(histogram) < 5 {
		return false
	}
	D := histogram[len(histogram)-2]
	E := histogram[len(histogram)-1]

	if E < 0 {
		return true
	}
	if D < 0 {
		return true
	}
	return false
}

// 判断DEA趋势
func IsDEAUP(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	if len(closePrices) < slowPeriod+signalPeriod+1 {
		return false
	}
	_, DEA, histogram := CalculateMACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	if len(histogram) < 5 {
		return false
	}
	return DEA[len(DEA)-1] > 0
}

// 判断DEA趋势
func IsDEADOWN(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	if len(closePrices) < slowPeriod+signalPeriod+1 {
		return false
	}
	_, DEA, histogram := CalculateMACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	if len(histogram) < 5 {
		return false
	}
	return DEA[len(DEA)-1] < 0
}

// 判断DIF趋势
func IsDIFUP(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	if len(closePrices) < slowPeriod+signalPeriod+1 {
		return false
	}
	DIF, _, histogram := CalculateMACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	if len(histogram) < 5 {
		return false
	}
	return DIF[len(DIF)-1] > 0
}

// 判断DEA趋势
func IsDIFDOWN(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	if len(closePrices) < slowPeriod+signalPeriod+1 {
		return false
	}
	DIF, _, histogram := CalculateMACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	if len(histogram) < 5 {
		return false
	}
	return DIF[len(DIF)-1] < 0
}

//为正
func UPUP(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	if len(closePrices) < slowPeriod+signalPeriod+1 {
		return false
	}

	_, _, histogram := CalculateMACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	if len(histogram) < 3 {
		return false
	}

	C := histogram[len(histogram)-3]
	D := histogram[len(histogram)-2]
	E := histogram[len(histogram)-1]

	return E > D || D > C
}

//为负
func DownDown(closePrices []float64, fastPeriod, slowPeriod, signalPeriod int) bool {
	if len(closePrices) < slowPeriod+signalPeriod+1 {
		return false
	}

	_, _, histogram := CalculateMACD(closePrices, fastPeriod, slowPeriod, signalPeriod)
	if len(histogram) < 3 {
		return false
	}

	C := histogram[len(histogram)-3]
	D := histogram[len(histogram)-2]
	E := histogram[len(histogram)-1]

	return E < D || D < C
}
