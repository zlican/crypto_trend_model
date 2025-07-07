# 加密货币趋势监控工具

这是一个用于监控 BTC 和 ETH 趋势的工具，基于 MA25、EMA144 和 EMA169 指标进行趋势判断。

## 功能特点

- 支持监控 BTC 和 ETH 的价格趋势
- 支持 1 小时、4 小时和 1 天三种时间周期
- 基于 MA25、EMA144 和 EMA169 指标进行趋势判断
- 趋势判断规则：
  - 上涨趋势：当前收盘价 > EMA25 且 > EMA144
  - 下跌趋势：当前收盘价 < EMA25 且 < EMA169
  - 震荡行情：其他情况
- 定时自动监控，默认每小时执行一次
- 支持控制台输出和文件日志记录
- 提供 HTTP API 接口，方便外部程序获取趋势数据
- 集成 Rainmeter 组件，可在桌面实时显示趋势状态

## 安装和使用

### 前置条件

- Go 1.16 或更高版本
- Rainmeter 4.0 或更高版本（可选，用于桌面显示）

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
./crypto_trend_monitor

# 或者使用Go命令运行
go run .
```

## 配置

配置参数在 `config/config.go` 文件中定义，可以根据需要修改：

- `APIBaseURL`: 币安 API 的基础 URL
- `KlineEndpoint`: K 线数据的 API 端点
- `Symbols`: 要监控的交易对列表
- `Intervals`: 要监控的时间周期列表
- `MA25Period`: MA25 的周期
- `EMA144Period`: EMA144 的周期
- `EMA169Period`: EMA169 的周期
- `MonitorInterval`: 监控频率（小时）
- `EnableAPIServer`: 是否启用 API 服务器
- `APIServerPort`: API 服务器端口

## API 接口

程序提供了以下 HTTP API 接口：

### 获取所有趋势数据

```
GET /api/trend?interval=1h
```

参数：

- `interval`: 时间周期，可选值：1h, 4h, 1d，默认为 1h

返回示例：

```json
{
  "btc": {
    "symbol": "BTC",
    "interval": "1h",
    "trend": "上涨趋势",
    "current_price": 60123.45,
    "ma25": 59876.32,
    "ema144": 58765.43,
    "ema169": 58234.56,
    "time": "2023-05-01 12:00:00"
  },
  "eth": {
    "symbol": "ETH",
    "interval": "1h",
    "trend": "震荡行情",
    "current_price": 3245.67,
    "ma25": 3210.45,
    "ema144": 3250.32,
    "ema169": 3180.76,
    "time": "2023-05-01 12:00:00"
  }
}
```

### 获取 BTC 趋势数据

```
GET /api/trend/btc?interval=1h&format=text
```

参数：

- `interval`: 时间周期，可选值：1h, 4h, 1d，默认为 1h
- `format`: 返回格式，可选值：json（默认）, text（适用于 Rainmeter）

JSON 格式返回示例：

```json
{
  "symbol": "BTC",
  "interval": "1h",
  "trend": "上涨趋势",
  "current_price": 60123.45,
  "ma25": 59876.32,
  "ema144": 58765.43,
  "ema169": 58234.56,
  "time": "2023-05-01 12:00:00"
}
```

文本格式返回示例：

```
BTC Trend: 上涨趋势
```

### 获取 ETH 趋势数据

```
GET /api/trend/eth?interval=1h&format=text
```

参数与 BTC 接口相同。

## Rainmeter 集成

本项目提供了 Rainmeter 皮肤，可以在桌面上实时显示 BTC 和 ETH 的趋势状态。

### 安装步骤

1. 安装 Rainmeter（如果尚未安装）
2. 复制`rainmeter/CryptoTrendMonitor.ini`文件到 Rainmeter 的皮肤目录
3. 在 Rainmeter 管理器中加载该皮肤

### 配置 Rainmeter 皮肤

在皮肤中，您可以修改以下参数：

- `UpdateInterval`: 更新频率（秒）
- `APIHost`: API 服务器主机名，默认为 localhost
- `APIPort`: API 服务器端口，默认为 8080
- `Interval`: 时间周期，可选值：1h, 4h, 1d

您还可以通过点击皮肤上的时间周期按钮（1H、4H、1D）来切换不同的时间周期。

## 日志

日志文件默认保存在 `logs` 目录下：

- `trend_analysis_YYYYMMDD.log`: 趋势分析结果日志
- `error_YYYYMMDD.log`: 错误日志

## 项目结构

- `main.go`: 主程序入口
- `config/config.go`: 配置参数
- `utils/binance_client.go`: 币安 API 客户端
- `utils/indicators.go`: 技术指标计算
- `utils/trend_analyzer.go`: 趋势分析
- `utils/output.go`: 输出和日志管理
- `utils/api_server.go`: API 服务器
- `rainmeter/CryptoTrendMonitor.ini`: Rainmeter 皮肤配置

## 许可证

[MIT License](LICENSE)
