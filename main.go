package main

import (
	"crypto_trend_monitor/config"
	"crypto_trend_monitor/model"
	"crypto_trend_monitor/utils"
	"database/sql"
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

var (
	db *sql.DB
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

		go func() {
			if err := apiServer.Start(); err != nil {
				log.Printf("API服务器启动失败: %v", err)
			}
		}()
	}

	model.InitDB()
	db = model.DB

	// ✅ 首次立即执行
	log.Printf("[TrendMonitor] 首次立即执行: %s", time.Now().Format("15:04:05"))
	results := runAnalysis(analyzer, output)
	if apiServer != nil && len(results) > 0 {
		apiServer.UpdateResults(results)
	}

	// ✅ 计算下一次 minute % 5 == 0 的时间
	now := time.Now()
	minutesToNext := 5 - (now.Minute() % 5)
	if minutesToNext == 0 {
		minutesToNext = 5
	}
	nextAligned := now.Truncate(time.Minute).Add(time.Duration(minutesToNext) * time.Minute)
	delay := time.Until(nextAligned)

	log.Printf("[TrendMonitor] 下一次对齐执行时间: %s（等待 %v）", nextAligned.Format("15:04:05"), delay)

	// ✅ 通道控制优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// ✅ 启动延迟后执行一次 + 开始固定周期执行
	go func() {
		time.Sleep(delay)

		log.Printf("[TrendMonitor] 对齐执行: %s", time.Now().Format("15:04:05"))
		results := runAnalysis(analyzer, output)
		if apiServer != nil && len(results) > 0 {
			apiServer.UpdateResults(results)
		}

		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Printf("[TrendMonitor] 周期触发: %s", time.Now().Format("15:04:05"))
				results := runAnalysis(analyzer, output)
				if apiServer != nil && len(results) > 0 {
					apiServer.UpdateResults(results)
				}
			case <-sigChan:
				log.Println("接收到退出信号，程序正在退出...")
				return
			}
		}
	}()

	// 阻塞主协程，直到收到退出信号
	<-sigChan
	log.Println("程序已退出。")
}

// runAnalysis 运行一次趋势分析
func runAnalysis(analyzer *utils.TrendAnalyzer, output *utils.OutputManager) []*utils.TrendResult {
	log.Println("开始执行趋势分析...")
	time.Sleep(7 * time.Second) //等待当前K线出来

	// 分析所有趋势
	results := analyzer.AnalyzeAllTrends(db)

	// 记录结果
	if err := output.LogTrendResults(results); err != nil {
		output.LogError(err)
	}

	log.Println("趋势分析完成")

	return results
}
