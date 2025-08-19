package gin_lin

import "time"

// 获取今天剩余时间
func GetRemainingSecondsToday() time.Duration {
	now := time.Now()

	// 获取今天的最后一秒（23:59:59.999999999）
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(),
		23, 59, 59, 999999999, now.Location())

	// 计算剩余时间并返回Duration
	return endOfDay.Sub(now)
}
