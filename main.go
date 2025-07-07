package main

import (
	"crypto_trend_monitor/config"
	"crypto_trend_monitor/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 定义常量
const (
	BinanceAPIBaseURL = "https://fapi.binance.com"
	KlineEndpoint     = "/fapi/v1/klines"
)

func main() {
	log.Println("开始运行币种趋势监控程序...")

	// 初始化输出管理器
	output := utils.NewOutputManager()
	if err := output.Init(); err != nil {
		log.Fatalf("初始化输出管理器失败: %v", err)
	}

	// 创建趋势分析器
	analyzer := utils.NewTrendAnalyzer()

	// 创建API服务器
	var apiServer *utils.TrendAPI
	if config.GlobalConfig.EnableAPIServer {
		apiServer = utils.NewTrendAPI(config.GlobalConfig.APIServerPort, analyzer)

		// 在goroutine中启动API服务器
		go func() {
			if err := apiServer.Start(); err != nil {
				log.Printf("API服务器启动失败: %v", err)
			}
		}()
	}

	// 首次运行
	results := runAnalysis(analyzer, output)

	// 如果启用了API服务器，更新结果
	if apiServer != nil && len(results) > 0 {
		apiServer.UpdateResults(results)
	}

	// 设置定时器，每小时执行一次
	ticker := time.NewTicker(time.Duration(config.GlobalConfig.MonitorInterval) * time.Minute)
	defer ticker.Stop()

	// 设置信号处理，以便优雅地退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("程序已启动，按 Ctrl+C 退出...")

	// 主循环
	for {
		select {
		case <-ticker.C:
			results := runAnalysis(analyzer, output)

			// 如果启用了API服务器，更新结果
			if apiServer != nil && len(results) > 0 {
				apiServer.UpdateResults(results)
			}
		case <-sigChan:
			log.Println("接收到退出信号，程序正在退出...")
			return
		}
	}
}

// runAnalysis 运行一次趋势分析
func runAnalysis(analyzer *utils.TrendAnalyzer, output *utils.OutputManager) []*utils.TrendResult {
	log.Println("开始执行趋势分析...")

	// 分析所有趋势
	results := analyzer.AnalyzeAllTrends()

	// 记录结果
	if err := output.LogTrendResults(results); err != nil {
		output.LogError(err)
	}

	log.Println("趋势分析完成")

	return results
}
