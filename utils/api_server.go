package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// TrendAPI API服务器，提供趋势数据接口
type TrendAPI struct {
	Port          int
	analyzer      *TrendAnalyzer
	latestResults map[string]*TrendResult // 按symbol存储最新结果
	mu            sync.RWMutex
}

// NewTrendAPI 创建新的API服务器
func NewTrendAPI(port int, analyzer *TrendAnalyzer) *TrendAPI {
	return &TrendAPI{
		Port:          port,
		analyzer:      analyzer,
		latestResults: make(map[string]*TrendResult),
	}
}

// Start 启动API服务器
func (api *TrendAPI) Start() error {

	http.HandleFunc("/api/trend/btc", api.handleTrendBTC)
	http.HandleFunc("/api/trend/eth", api.handleTrendETH)

	// 启动服务器
	addr := fmt.Sprintf(":%d", api.Port)
	log.Printf("API服务器启动在 http://localhost%s", addr)

	return http.ListenAndServe(addr, nil)
}

// UpdateResults 更新最新的趋势结果
func (api *TrendAPI) UpdateResults(results []*TrendResult) {
	api.mu.Lock()
	defer api.mu.Unlock()

	// 按symbol和interval组织结果
	for _, result := range results {
		key := fmt.Sprintf("%s_%s", result.Symbol, result.Interval)
		api.latestResults[key] = result
	}
}

// handleTrendBTC 处理获取BTC趋势的请求
func (api *TrendAPI) handleTrendBTC(w http.ResponseWriter, r *http.Request) {
	api.mu.RLock()
	defer api.mu.RUnlock()

	// 获取interval参数，默认为1h
	interval := r.URL.Query().Get("interval")
	if interval == "" {
		interval = "1h"
	}

	// 获取BTC趋势
	btcKey := fmt.Sprintf("BTCUSDT_%s", interval)
	if btcResult, ok := api.latestResults[btcKey]; ok {
		// 根据请求格式返回不同的响应
		apiStatus := "unknown"
		if btcResult.Status == UP {
			apiStatus = "UP"
		} else if btcResult.Status == DOWN {
			apiStatus = "DOWN"
		} else if btcResult.Status == UPEMA {
			apiStatus = "UPEMA"
		} else {
			apiStatus = "DOWNEMA"
		}

		if r.URL.Query().Get("format") == "text" {
			// 纯文本格式，适合Rainmeter
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "BTC Trend: %s", apiStatus)
		} else {
			// JSON格式
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"symbol":   "BTC",
				"interval": btcResult.Interval,
				"trend":    apiStatus,
				"ema25":    btcResult.EMA25,
				"ema50":    btcResult.EMA50,
				"time":     btcResult.Time.Format("2006-01-02 15:04:05"),
			})
		}
	} else {
		// 数据不可用
		if r.URL.Query().Get("format") == "text" {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "BTC Trend: unknown")
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "BTC Trend: unknown",
			})
		}
	}
}

// handleTrendETH 处理获取ETH趋势的请求
func (api *TrendAPI) handleTrendETH(w http.ResponseWriter, r *http.Request) {
	api.mu.RLock()
	defer api.mu.RUnlock()

	// 获取interval参数，默认为1h
	interval := r.URL.Query().Get("interval")
	if interval == "" {
		interval = "1h"
	}

	// 获取ETH趋势
	ethKey := fmt.Sprintf("ETHUSDT_%s", interval)
	if ethResult, ok := api.latestResults[ethKey]; ok {
		// 根据请求格式返回不同的响应
		apiStatus := "unknown"
		if ethResult.Status == UP {
			apiStatus = "UP"
		} else if ethResult.Status == DOWN {
			apiStatus = "DOWN"
		} else if ethResult.Status == UPEMA {
			apiStatus = "UPEMA"
		} else {
			apiStatus = "DOWNEMA"

		}

		if r.URL.Query().Get("format") == "text" {
			// 纯文本格式，适合Rainmeter
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "ETH Trend: %s", apiStatus)
		} else {
			// JSON格式
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"symbol":   "ETH",
				"interval": ethResult.Interval,
				"trend":    apiStatus,
				"ema25":    ethResult.EMA25,
				"ema50":    ethResult.EMA50,
				"time":     ethResult.Time.Format("2006-01-02 15:04:05"),
			})
		}
	} else {
		// 数据不可用
		if r.URL.Query().Get("format") == "text" {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "ETH Trend: unknown")
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "ETH Trend: unknown",
			})
		}
	}
}
