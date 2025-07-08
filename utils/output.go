package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// OutputManager 输出管理器
type OutputManager struct {
	LogDir     string
	ConsoleLog bool
	FileLog    bool
}

// NewOutputManager 创建输出管理器
func NewOutputManager() *OutputManager {
	return &OutputManager{
		LogDir:     "logs",
		ConsoleLog: true,
		FileLog:    true,
	}
}

// Init 初始化输出管理器
func (o *OutputManager) Init() error {
	if o.FileLog {
		// 创建日志目录
		if err := os.MkdirAll(o.LogDir, 0755); err != nil {
			return fmt.Errorf("创建日志目录失败: %v", err)
		}
	}
	return nil
}

// LogTrendResults 记录趋势分析结果
func (o *OutputManager) LogTrendResults(results []*TrendResult) error {
	// 控制台输出
	if o.ConsoleLog {
		fmt.Println("===== 趋势分析结果 =====")
		fmt.Printf("分析时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Println("------------------------")

		for _, result := range results {
			fmt.Printf("%s %s: %s\n", result.Symbol, result.Interval, result.Status)
			fmt.Printf("  当前价格: %.2f\n", result.CurrentPrice)
			fmt.Printf("  EMA25: %.2f\n", result.EMA25)
			fmt.Printf("  EMA50: %.2f\n", result.EMA50)
			fmt.Println("------------------------")
		}
	}

	// 文件输出
	if o.FileLog {
		// 创建日志文件
		logFileName := fmt.Sprintf("trend_analysis_%s.log", time.Now().Format("20060102"))
		logFilePath := filepath.Join(o.LogDir, logFileName)

		// 追加模式打开文件
		f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("打开日志文件失败: %v", err)
		}
		defer f.Close()

		// 写入日志头
		f.WriteString(fmt.Sprintf("===== 趋势分析结果 (%s) =====\n", time.Now().Format("2006-01-02 15:04:05")))

		// 写入每个结果
		for _, result := range results {
			f.WriteString(FormatTrendResult(result) + "\n")
		}

		f.WriteString("==============================\n\n")
	}

	return nil
}

// LogError 记录错误信息
func (o *OutputManager) LogError(err error) {
	if o.ConsoleLog {
		log.Printf("错误: %v", err)
	}

	if o.FileLog {
		// 创建错误日志文件
		logFileName := fmt.Sprintf("error_%s.log", time.Now().Format("20060102"))
		logFilePath := filepath.Join(o.LogDir, logFileName)

		// 追加模式打开文件
		f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("打开错误日志文件失败: %v", err)
			return
		}
		defer f.Close()

		// 写入错误信息
		f.WriteString(fmt.Sprintf("[%s] %v\n", time.Now().Format("2006-01-02 15:04:05"), err))
	}
}
