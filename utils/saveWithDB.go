package utils

import (
	"database/sql"
	"fmt"
	"time"
)

// 保存趋势结果
func SaveTrendResult(db *sql.DB, result *TrendResult) error {
	// 根据 interval 选择表名
	var tableName string
	switch result.Interval {
	case "5m":
		tableName = "symbol_5m"
	case "15m":
		tableName = "symbol_15m"
	case "1h":
		tableName = "symbol_1h"
	case "4h":
		tableName = "symbol_4h"
	case "1d":
		tableName = "symbol_1d"
	case "3d":
		tableName = "symbol_3d"
	default:
		return fmt.Errorf("不支持的 interval: %s", result.Interval)
	}

	// timestamp 这里你可以用 K线最后一根的时间戳，
	// 如果你现在没有，那就用当前时间（秒级）
	timestamp := time.Now().Unix()

	// SQL：插入或更新
	query := fmt.Sprintf(`
		INSERT INTO %s (symbol, timestamp, status)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			status = VALUES(status),
			updated_at = CURRENT_TIMESTAMP
	`, tableName)

	_, err := db.Exec(query, result.Symbol, timestamp, result.Status)
	if err != nil {
		return fmt.Errorf("保存到数据库失败: %v", err)
	}

	return nil
}
