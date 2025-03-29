package utils

import (
	"math"
	"math/rand"
	"time"
)

// 生成随机整数
func Random(n int32) int32 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(int32(math.Pow10(int(n))))
}
