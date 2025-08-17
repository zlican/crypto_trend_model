package utils

import (
	"crypto_trend_monitor/config"
	"database/sql"
	"fmt"
	"time"
)

// TrendStatus 表示趋势状态
type TrendStatus string

const (
	RANGE    TrendStatus = "RANGE"
	BUYMACD  TrendStatus = "BUYMACD"
	SELLMACD TrendStatus = "SELLMACD"
	UP       TrendStatus = "UP"
	DOWN     TrendStatus = "DOWN"
)

// TrendResult 趋势分析结果
type TrendResult struct {
	Symbol   string
	Interval string
	Status   TrendStatus
	EMA25    float64
	EMA50    float64
	Time     time.Time
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
func (a *TrendAnalyzer) AnalyzeTrend(symbol, interval string, db *sql.DB) (*TrendResult, error) {
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
	price := closePrices[len(closePrices)-1]
	ema25 := a.indicators["EMA25"].Calculate(closePrices)
	ema50 := a.indicators["EMA50"].Calculate(closePrices)
	ma60 := CalculateMA(closePrices, 60)
	UpMACD := IsAboutToGoldenCross(closePrices, 6, 13, 5)
	DownMACD := IsAboutToDeadCross(closePrices, 6, 13, 5)
	XUpMACD := IsGolden(closePrices, 6, 13, 5)
	XDownMACD := IsDead(closePrices, 6, 13, 5)

	var BuyMACD, SellMACD bool
	UPEMA := ema25 > ema50
	DOWNEMA := ema25 < ema50
	if UPEMA && UpMACD && price > ema25 && (price > ma60 || ma60 < ema25) { //金叉回调
		BuyMACD = true
	} else if DOWNEMA && XUpMACD && price > ema25 && (price > ma60 || ma60 < ema25) { //死叉反转
		BuyMACD = true
	} else if DOWNEMA && DownMACD && price < ema25 && (price < ma60 || ma60 > ema25) {
		SellMACD = true
	} else if UPEMA && XDownMACD && price < ema25 && (price < ma60 || ma60 > ema25) {
		SellMACD = true
	} else {
		BuyMACD = false
		SellMACD = false
	}

	// 判断趋势
	var status TrendStatus

	if BuyMACD {
		status = BUYMACD
	} else if SellMACD {
		status = SELLMACD
	} else {
		status = RANGE
	}

	res := &TrendResult{
		Symbol:   symbol,
		Interval: interval,
		Status:   status,
		EMA25:    ema25,
		EMA50:    ema50,
		Time:     time.Now(),
	}

	if err := SaveTrendResult(db, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AnalyzeAllTrends 分析所有配置的币种和时间周期的趋势
func (a *TrendAnalyzer) AnalyzeAllTrends(db *sql.DB) []*TrendResult {
	results := make([]*TrendResult, 0)

	for _, symbol := range config.GlobalConfig.Symbols {
		for _, interval := range config.GlobalConfig.Intervals {
			result, err := a.AnalyzeTrend(symbol, interval, db)
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
		"[%s] %s %s: 当前价格=%.2f, EMA25=%.2f, EMA50=%.2f, 趋势=%s",
		result.Time.Format("2006-01-02 15:04:05"),
		result.Symbol,
		result.Interval,
		result.EMA25,
		result.EMA50,
		result.Status,
	)
}
