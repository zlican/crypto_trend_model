package config

// Config 包含程序的配置参数
type Config struct {
	// API配置
	APIBaseURL    string
	KlineEndpoint string

	// 监控的交易对
	Symbols []string

	// 监控的时间周期
	Intervals []string

	//代理
	ProxyURL string

	// 技术指标参数
	EMA25Period  int
	EMA50Period  int
	EMA120Period int

	// 监控频率（分钟）
	MonitorInterval int

	// API服务器配置
	EnableAPIServer bool
	APIServerPort   int
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		APIBaseURL:      "https://fapi.binance.com",
		KlineEndpoint:   "/fapi/v1/klines",
		Symbols:         []string{"BTCUSDT", "ETHUSDT"},
		Intervals:       []string{"5m", "15m", "1h", "4h", "1d", "3d"},
		ProxyURL:        "http://127.0.0.1:10809",
		EMA25Period:     25,
		EMA50Period:     50,
		MonitorInterval: 15, // 每15分钟
		EnableAPIServer: true,
		APIServerPort:   8080,
	}
}

// 全局配置实例
var GlobalConfig = DefaultConfig()
