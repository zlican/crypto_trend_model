package utils

import "crypto_trend_monitor/config"

// Indicator 接口定义了技术指标的通用行为
type Indicator interface {
	Calculate(prices []float64) float64
	GetPeriod() int
}

// SimpleMA 简单移动平均线
type SimpleMA struct {
	Period int
}

// Calculate 计算简单移动平均线
func (ma *SimpleMA) Calculate(prices []float64) float64 {
	if len(prices) < ma.Period {
		return 0
	}

	sum := 0.0
	for i := len(prices) - ma.Period; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(ma.Period)
}

// GetPeriod 返回周期
func (ma *SimpleMA) GetPeriod() int {
	return ma.Period
}

// EMA 指数移动平均线
type EMA struct {
	Period int
}

// Calculate 计算指数移动平均线
func (ema *EMA) Calculate(prices []float64) float64 {
	if len(prices) < ema.Period {
		return 0 // 数据不足，无法计算
	}

	//prices = ReverseSlice(prices)

	alpha := 2.0 / float64(ema.Period+1)
	emas := make([]float64, len(prices))

	// 计算初始 EMA 值，使用前 Period 个价格的简单平均
	sum := 0.0
	for i := 0; i < ema.Period; i++ {
		sum += prices[i]
	}
	emas[ema.Period-1] = sum / float64(ema.Period)

	// 迭代计算 EMA
	for i := ema.Period; i < len(prices); i++ {
		emas[i] = (prices[i]-emas[i-1])*alpha + emas[i-1]
	}

	return emas[len(prices)-1]
}

// GetPeriod 返回周期
func (ema *EMA) GetPeriod() int {
	return ema.Period
}

// NewIndicators 创建所有需要的技术指标
func NewIndicators() map[string]Indicator {
	return map[string]Indicator{
		"EMA25":  &EMA{Period: config.GlobalConfig.EMA25Period},
		"EMA144": &EMA{Period: config.GlobalConfig.EMA144Period},
		"EMA169": &EMA{Period: config.GlobalConfig.EMA169Period},
	}
}

// GetMaxPeriod 获取所有指标中最大的周期值
func GetMaxPeriod(indicators map[string]Indicator) int {
	maxPeriod := 0
	for _, indicator := range indicators {
		if indicator.GetPeriod() > maxPeriod {
			maxPeriod = indicator.GetPeriod()
		}
	}
	return maxPeriod
}

//反转prices
func ReverseSlice(s []float64) []float64 {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
