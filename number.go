package gin_lin

import (
	"math"
	"math/rand"
	"time"
)

/**
 * @description: 生成指定位数随机整数
 * @param {int32} n 位数
 */
func Random(n int32) int32 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(int32(math.Pow10(int(n))))
}

/**
 * @description: 生成随机数，指定范围
 * @param {int} min 最小值，包含
 * @param {int} max 最大值，不包含
 */
func RandomInt(min int, max int) int {
	// 生成[min, max)范围的随机数
	return rand.Intn(max-min) + min
}

/**
 * @description: 生成随机数，指定范围
 * @param {int} min 最小值，包含
 * @param {int} max 最大值，不包含
 */
func RandomFloat(start float64, end float64) float64 {
	return start + rand.Float64()*(end-start)
}

/**
 * @description: 保留两位浮点数小数
 */
func TruncateToTwoDecimal(num float64) float64 {
	return math.Trunc(num*100) / 100
}
