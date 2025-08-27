package util

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	// 用于生成唯一序列号的计数器
	counter uint64
	// 保护计数器的互斥锁
	counterMutex sync.Mutex
	// 随机数生成器实例
	rng *rand.Rand
)

// GenerateOrderID 生成基于日期和序列号的订单ID
// 格式: YYYYMMDDHHMMSS + 6位序列号
// 例如: 20240101120000000001
func GenerateOrderID() string {
	now := time.Now()
	// 格式化时间为 YYYYMMDDHHMMSS
	timeStr := now.Format("20060102150405")

	// 使用互斥锁保护计数器
	counterMutex.Lock()
	counter++
	seq := counter % 1000000 // 保持6位数字
	if seq == 0 {
		seq = 1 // 避免000000
	}
	counterMutex.Unlock()

	return fmt.Sprintf("%s%06d", timeStr, seq)
}

// GenerateOrderIDWithPrefix 生成带前缀的订单ID
// 格式: 前缀 + YYYYMMDDHHMMSS + 6位随机数
// 例如: ORD20240806161530123456
func GenerateOrderIDWithPrefix(prefix string) string {
	return prefix + GenerateOrderID()
}

// GenerateShortOrderID 生成短格式订单ID
// 格式: YYMMDDHHMMSS + 4位随机数
// 例如: 240806161530123456
func GenerateShortOrderID() string {
	// 获取当前时间
	now := time.Now()

	// 格式化日期时间部分: YYMMDDHHMMSS
	dateTimePart := now.Format("060102150405")

	// 生成4位随机数
	randomPart := rng.Intn(10000) // 0-9999

	// 组合成完整的订单ID
	orderID := fmt.Sprintf("%s%04d", dateTimePart, randomPart)

	return orderID
}

// init 初始化随机数生成器
func init() {
	// 使用当前时间纳秒作为种子创建新的随机数生成器
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}
