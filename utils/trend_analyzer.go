package utils

import (
	"crypto_trend_monitor/config"
	"fmt"
	"time"
)

// TrendStatus 表示趋势状态
type TrendStatus string

const (
	UpTrend    TrendStatus = "上涨趋势"
	DownTrend  TrendStatus = "下跌趋势"
	RangeBound TrendStatus = "震荡行情"
)

// TrendResult 趋势分析结果
type TrendResult struct {
	Symbol       string
	Interval     string
	Status       TrendStatus
	CurrentPrice float64
	EMA25        float64
	EMA144       float64
	EMA169       float64
	Time         time.Time
}

// TrendAnalyzer 趋势分析器
type TrendAnalyzer struct {
	client     *BinanceClient
	indicators map[string]Indicator
}

// NewTrendAnalyzer 创建趋势分析器
func NewTrendAnalyzer() *TrendAnalyzer {
	return &TrendAnalyzer{
		client:     NewBinanceClient(),
		indicators: NewIndicators(),
	}
}

// AnalyzeTrend 分析特定币种和时间周期的趋势
func (a *TrendAnalyzer) AnalyzeTrend(symbol, interval string) (*TrendResult, error) {
	// 获取最大周期值，并多获取一些数据以确保足够
	maxPeriod := GetMaxPeriod(a.indicators)
	limit := 499

	// 获取K线数据
	klines, err := a.client.GetKlines(symbol, interval, limit)
	if err != nil {
		return nil, fmt.Errorf("获取K线数据失败: %v", err)
	}

	if len(klines) < maxPeriod {
		return nil, fmt.Errorf("获取的K线数据不足: %d/%d", len(klines), maxPeriod)
	}

	// 提取收盘价
	closePrices := ExtractClosePrices(klines)

	// 计算指标
	currentPrice := closePrices[len(closePrices)-1]
	ema25 := a.indicators["EMA25"].Calculate(closePrices)
	ema144 := a.indicators["EMA144"].Calculate(closePrices)
	ema169 := a.indicators["EMA169"].Calculate(closePrices)

	// 判断趋势
	var status TrendStatus
	if currentPrice > ema25 && currentPrice > ema144 {
		status = UpTrend
	} else if currentPrice < ema25 && currentPrice < ema169 {
		status = DownTrend
	} else {
		status = RangeBound
	}

	return &TrendResult{
		Symbol:       symbol,
		Interval:     interval,
		Status:       status,
		CurrentPrice: currentPrice,
		EMA25:        ema25,
		EMA144:       ema144,
		EMA169:       ema169,
		Time:         time.Now(),
	}, nil
}

// AnalyzeAllTrends 分析所有配置的币种和时间周期的趋势
func (a *TrendAnalyzer) AnalyzeAllTrends() []*TrendResult {
	results := make([]*TrendResult, 0)

	for _, symbol := range config.GlobalConfig.Symbols {
		for _, interval := range config.GlobalConfig.Intervals {
			result, err := a.AnalyzeTrend(symbol, interval)
			if err != nil {
				fmt.Printf("分析 %s %s 趋势失败: %v\n", symbol, interval, err)
				continue
			}
			results = append(results, result)
		}
	}

	return results
}

// FormatTrendResult 格式化趋势结果为字符串
func FormatTrendResult(result *TrendResult) string {
	return fmt.Sprintf(
		"[%s] %s %s: 当前价格=%.2f, EMA25=%.2f, EMA144=%.2f, EMA169=%.2f, 趋势=%s",
		result.Time.Format("2006-01-02 15:04:05"),
		result.Symbol,
		result.Interval,
		result.CurrentPrice,
		result.EMA25,
		result.EMA144,
		result.EMA169,
		result.Status,
	)
}
