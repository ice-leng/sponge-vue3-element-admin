package util

import (
	"regexp"
	"testing"
	"time"
)

// TestGenerateOrderID 测试生成订单ID
func TestGenerateOrderID(t *testing.T) {
	// 生成订单ID
	orderID := GenerateOrderID()

	// 检查长度 (14位日期时间 + 6位随机数 = 20位)
	if len(orderID) != 20 {
		t.Errorf("Expected order ID length 20, got %d", len(orderID))
	}

	// 检查格式 (应该全是数字)
	matched, _ := regexp.MatchString(`^\d{20}$`, orderID)
	if !matched {
		t.Errorf("Order ID should contain only digits, got %s", orderID)
	}

	// 检查日期部分是否合理 (前14位应该是当前时间附近)
	now := time.Now()
	expectedPrefix := now.Format("20060102150405")
	actualPrefix := orderID[:14]

	// 允许1秒的时间差
	expectedTime, _ := time.Parse("20060102150405", expectedPrefix)
	actualTime, _ := time.Parse("20060102150405", actualPrefix)
	diff := actualTime.Sub(expectedTime)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("Order ID timestamp seems incorrect. Expected around %s, got %s", expectedPrefix, actualPrefix)
	}

	t.Logf("Generated order ID: %s", orderID)
}

// TestGenerateOrderIDWithPrefix 测试生成带前缀的订单ID
func TestGenerateOrderIDWithPrefix(t *testing.T) {
	prefix := "ORD"
	orderID := GenerateOrderIDWithPrefix(prefix)

	// 检查长度 (3位前缀 + 20位ID = 23位)
	if len(orderID) != 23 {
		t.Errorf("Expected order ID length 23, got %d", len(orderID))
	}

	// 检查前缀
	if orderID[:3] != prefix {
		t.Errorf("Expected prefix %s, got %s", prefix, orderID[:3])
	}

	// 检查数字部分格式
	numberPart := orderID[3:]
	matched, _ := regexp.MatchString(`^\d{20}$`, numberPart)
	if !matched {
		t.Errorf("Order ID number part should contain only digits, got %s", numberPart)
	}

	t.Logf("Generated order ID with prefix: %s", orderID)
}

// TestGenerateShortOrderID 测试生成短格式订单ID
func TestGenerateShortOrderID(t *testing.T) {
	// 生成短订单ID
	orderID := GenerateShortOrderID()

	// 检查长度 (12位日期时间 + 4位随机数 = 16位)
	if len(orderID) != 16 {
		t.Errorf("Expected short order ID length 16, got %d", len(orderID))
	}

	// 检查格式 (应该全是数字)
	matched, _ := regexp.MatchString(`^\d{16}$`, orderID)
	if !matched {
		t.Errorf("Short order ID should contain only digits, got %s", orderID)
	}

	t.Logf("Generated short order ID: %s", orderID)
}

// TestOrderIDUniqueness 测试订单ID的唯一性
func TestOrderIDUniqueness(t *testing.T) {
	ids := make(map[string]bool)
	duplicateCount := 0
	totalCount := 1000

	// 生成1000个ID，检查重复率
	for i := 0; i < totalCount; i++ {
		id := GenerateOrderID()
		if ids[id] {
			duplicateCount++
			t.Logf("Duplicate ID found: %s", id)
		}
		ids[id] = true
	}

	// 计算重复率
	duplicateRate := float64(duplicateCount) / float64(totalCount) * 100
	t.Logf("Generated %d IDs, found %d duplicates (%.2f%%)", totalCount, duplicateCount, duplicateRate)

	// 重复率应该很低 (小于1%)
	if duplicateRate > 1.0 {
		t.Errorf("Duplicate rate too high: %.2f%%, expected < 1%%", duplicateRate)
	}
}

// BenchmarkGenerateOrderID 性能测试
func BenchmarkGenerateOrderID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateOrderID()
	}
}

// BenchmarkGenerateOrderIDWithPrefix 带前缀的性能测试
func BenchmarkGenerateOrderIDWithPrefix(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateOrderIDWithPrefix("ORD")
	}
}

// BenchmarkGenerateShortOrderID 短格式的性能测试
func BenchmarkGenerateShortOrderID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateShortOrderID()
	}
}
