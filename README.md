# 加密货币趋势监控工具

这是一个用于监控 BTC 和 ETH 趋势的工具，基于 MA25、EMA144 和 EMA169 指标进行趋势判断。

## 功能特点

- 支持监控 BTC 和 ETH 的价格趋势
- 支持 1 小时、4 小时和 1 天三种时间周期
- 基于 MA25、EMA144 和 EMA169 指标进行趋势判断
- 趋势判断规则：
  - 上涨趋势：当前收盘价 > MA25 且 > EMA144
  - 下跌趋势：当前收盘价 < MA25 且 < EMA169
  - 震荡行情：其他情况
- 定时自动监控，默认每小时执行一次
- 支持控制台输出和文件日志记录

## 安装和使用

### 前置条件

- Go 1.16 或更高版本

### 安装

```bash
# 克隆仓库
git clone [仓库URL]
cd [仓库目录]

# 构建项目
go build
```

### 运行

```bash
# 直接运行
./trend_monitor

# 或者使用Go命令运行
go run .
```

## 配置

配置参数在 `config.go` 文件中定义，可以根据需要修改：

- `APIBaseURL`: 币安 API 的基础 URL
- `KlineEndpoint`: K 线数据的 API 端点
- `Symbols`: 要监控的交易对列表
- `Intervals`: 要监控的时间周期列表
- `MA25Period`: MA25 的周期
- `EMA144Period`: EMA144 的周期
- `EMA169Period`: EMA169 的周期
- `MonitorInterval`: 监控频率（小时）

## 日志

日志文件默认保存在 `logs` 目录下：

- `trend_analysis_YYYYMMDD.log`: 趋势分析结果日志
- `error_YYYYMMDD.log`: 错误日志

## 项目结构

- `main.go`: 主程序入口
- `config.go`: 配置参数
- `binance_client.go`: 币安 API 客户端
- `indicators.go`: 技术指标计算
- `trend_analyzer.go`: 趋势分析
- `output.go`: 输出和日志管理

## 许可证

[MIT License](LICENSE)
