package utils

import (
	"crypto_trend_monitor/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// BinanceClient 币安API客户端
type BinanceClient struct {
	BaseURL    string
	HTTPClient *http.Client
	ProxyURL   string
}

// KlineData K线数据结构
type KlineData struct {
	OpenTime                 int64
	Open                     float64
	High                     float64
	Low                      float64
	Close                    float64
	Volume                   float64
	CloseTime                int64
	QuoteAssetVolume         float64
	NumberOfTrades           int64
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
}

// NewBinanceClient 创建一个新的币安客户端
func NewBinanceClient() *BinanceClient {
	return &BinanceClient{
		BaseURL: config.GlobalConfig.APIBaseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		ProxyURL: config.GlobalConfig.ProxyURL,
	}
}

// GetKlines 获取K线数据
func (c *BinanceClient) GetKlines(symbol, interval string, limit int) ([]KlineData, error) {
	urls := fmt.Sprintf("%s%s?symbol=%s&interval=%s&limit=%d",
		c.BaseURL, config.GlobalConfig.KlineEndpoint, symbol, interval, limit)

	proxyURL, _ := url.Parse(c.ProxyURL)
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   c.HTTPClient.Timeout,
	}

	var resp *http.Response
	var err error
	maxRetries := 3
	retryCount := 0
	retryDelay := 2 * time.Second

	for retryCount < maxRetries {
		resp, err = client.Get(urls)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		if resp != nil {
			resp.Body.Close()
		}

		retryCount++
		if retryCount >= maxRetries {
			if err != nil {
				return nil, fmt.Errorf("请求K线数据失败(已重试%d次): %v", maxRetries, err)
			}
			return nil, fmt.Errorf("API返回错误状态码(已重试%d次): %d", maxRetries, resp.StatusCode)
		}

		fmt.Printf("请求失败，%d秒后进行第%d次重试...\n", retryDelay/time.Second, retryCount+1)
		time.Sleep(retryDelay)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %v", err)
	}

	var rawKlines [][]interface{}
	if err := json.Unmarshal(body, &rawKlines); err != nil {
		return nil, fmt.Errorf("解析K线数据失败: %v", err)
	}

	klines := make([]KlineData, 0, len(rawKlines))
	for _, k := range rawKlines {
		if len(k) < 11 {
			return nil, fmt.Errorf("K线数据格式错误")
		}

		openTime, _ := k[0].(float64)
		open, _ := strconv.ParseFloat(k[1].(string), 64)
		high, _ := strconv.ParseFloat(k[2].(string), 64)
		low, _ := strconv.ParseFloat(k[3].(string), 64)
		close, _ := strconv.ParseFloat(k[4].(string), 64)
		volume, _ := strconv.ParseFloat(k[5].(string), 64)
		closeTime, _ := k[6].(float64)
		quoteVolume, _ := strconv.ParseFloat(k[7].(string), 64)
		numTrades, _ := k[8].(float64)
		takerBuyBaseVolume, _ := strconv.ParseFloat(k[9].(string), 64)
		takerBuyQuoteVolume, _ := strconv.ParseFloat(k[10].(string), 64)

		klines = append(klines, KlineData{
			OpenTime:                 int64(openTime),
			Open:                     open,
			High:                     high,
			Low:                      low,
			Close:                    close,
			Volume:                   volume,
			CloseTime:                int64(closeTime),
			QuoteAssetVolume:         quoteVolume,
			NumberOfTrades:           int64(numTrades),
			TakerBuyBaseAssetVolume:  takerBuyBaseVolume,
			TakerBuyQuoteAssetVolume: takerBuyQuoteVolume,
		})
	}

	return klines, nil
}

// ExtractClosePrices 从K线数据中提取收盘价
func ExtractClosePrices(klines []KlineData) []float64 {
	prices := make([]float64, len(klines))
	for i, kline := range klines {
		prices[i] = kline.Close
	}
	return prices
}

// ExtractClosePrices 从K线数据中提取收盘价
func ExtractOpensPrices(klines []KlineData) []float64 {
	prices := make([]float64, len(klines))
	for i, kline := range klines {
		prices[i] = kline.Open
	}
	return prices
}
