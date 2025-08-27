package datex

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestDatex 测试日期处理功能
func TestDatex(t *testing.T) {
	// 展示使用当前时间的情况
	showCurrentDateExample()

	// 展示使用指定日期的情况
	showSpecificDateExample()

	// 展示所有日期类型的示例
	showAllDateTypesExample()

	// 展示所有日期范围
	showDaysBetween()
}

func formatTime(t time.Time, format string) string {
	// 如果未指定格式，则使用标准的时间格式
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	// 使用指定的时间格式
	return t.Format(format)
}

// showCurrentDateExample 展示使用当前时间的示例
func showCurrentDateExample() {
	fmt.Println("\n===== 使用当前时间的示例 =====")

	// 使用当前时间
	now := time.Now()

	// 获取今天的开始和结束时间
	startTime, err := GetStartTime(DateTypeDay, now)
	if err != nil {
		fmt.Printf("获取开始时间失败: %v\n", err)
		return
	}

	endTime, err := GetEndTime(DateTypeDay, now)
	if err != nil {
		fmt.Printf("获取结束时间失败: %v\n", err)
		return
	}

	// 格式化并打印时间
	fmt.Printf("今天的开始时间: %s\n", formatTime(startTime, ""))
	fmt.Printf("今天的结束时间: %s\n", formatTime(endTime, ""))

	// 判断是否为今天
	isToday, err := IsToday(DateTypeDay, now)
	if err != nil {
		fmt.Printf("判断是否为今天失败: %v\n", err)
		return
	}
	fmt.Printf("是否为今天: %v\n", isToday)

	// 获取昨天的开始和结束时间
	// 使用前一天的日期
	yesterday := now.AddDate(0, 0, -1)
	prevStartTime, err := GetStartTime(DateTypeDay, yesterday)
	if err != nil {
		fmt.Printf("获取上一周期开始时间失败: %v\n", err)
		return
	}

	prevEndTime, err := GetEndTime(DateTypeDay, yesterday)
	if err != nil {
		fmt.Printf("获取上一周期结束时间失败: %v\n", err)
		return
	}

	// 格式化并打印时间
	fmt.Printf("昨天的开始时间: %s\n", formatTime(prevStartTime, ""))
	fmt.Printf("昨天的结束时间: %s\n", formatTime(prevEndTime, ""))
}

// showSpecificDateExample 展示使用指定日期的示例
func showSpecificDateExample() {
	fmt.Println("\n===== 使用指定日期的示例 =====")

	// 使用指定日期：2023年5月15日
	specificDate := time.Date(2023, 5, 15, 12, 30, 0, 0, time.Local)
	fmt.Printf("指定日期: %s\n", formatTime(specificDate, ""))

	// 获取该月的开始和结束时间
	startTime, err := GetStartTime(DateTypeMonth, specificDate)
	if err != nil {
		fmt.Printf("获取开始时间失败: %v\n", err)
		return
	}

	endTime, err := GetEndTime(DateTypeMonth, specificDate)
	if err != nil {
		fmt.Printf("获取结束时间失败: %v\n", err)
		return
	}

	// 格式化并打印时间
	fmt.Printf("2023年5月的开始时间: %s\n", formatTime(startTime, ""))
	fmt.Printf("2023年5月的结束时间: %s\n", formatTime(endTime, ""))

	// 获取上个月的开始和结束时间
	// 使用前一个月的日期
	prevMonth := specificDate.AddDate(0, -1, 0)
	prevStartTime, err := GetStartTime(DateTypeMonth, prevMonth)
	if err != nil {
		fmt.Printf("获取上一周期开始时间失败: %v\n", err)
		return
	}

	prevEndTime, err := GetEndTime(DateTypeMonth, prevMonth)
	if err != nil {
		fmt.Printf("获取上一周期结束时间失败: %v\n", err)
		return
	}

	// 格式化并打印时间
	fmt.Printf("2023年4月的开始时间: %s\n", formatTime(prevStartTime, ""))
	fmt.Printf("2023年4月的结束时间: %s\n", formatTime(prevEndTime, ""))
}

// showAllDateTypesExample 展示所有日期类型的示例
func showAllDateTypesExample() {
	fmt.Println("\n===== 所有日期类型的示例 =====")

	// 使用固定的测试日期：2023年5月15日（星期一）
	testDate := time.Now()
	fmt.Printf("测试日期: %s\n\n", formatTime(testDate, ""))

	// 展示日类型
	showDateTypeExample(testDate, DateTypeDay, "日")

	// 展示周类型
	showDateTypeExample(testDate, DateTypeWeek, "周")

	// 展示月类型
	showDateTypeExample(testDate, DateTypeMonth, "月")

	// 展示季度类型
	showDateTypeExample(testDate, DateTypeQuarter, "季度")

	// 展示年类型
	showDateTypeExample(testDate, DateTypeYear, "年")
}

// showDateTypeExample 展示特定日期类型的示例
func showDateTypeExample(date time.Time, dateType int, typeName string) {
	fmt.Printf("--- %s类型 ---\n", typeName)

	// 获取当前周期的开始和结束时间
	startTime, err := GetStartTime(dateType, date)
	if err != nil {
		fmt.Printf("获取开始时间失败: %v\n", err)
		return
	}

	endTime, err := GetEndTime(dateType, date)
	if err != nil {
		fmt.Printf("获取结束时间失败: %v\n", err)
		return
	}

	// 格式化并打印时间
	fmt.Printf("当前%s的开始时间: %s\n", typeName, formatTime(startTime, ""))
	fmt.Printf("当前%s的结束时间: %s\n", typeName, formatTime(endTime, ""))

	// 获取上一个周期的开始和结束时间
	// 根据日期类型计算上一个周期的日期
	var prevDate time.Time
	switch dateType {
	case DateTypeDay:
		prevDate = date.AddDate(0, 0, -1)
	case DateTypeWeek:
		prevDate = date.AddDate(0, 0, -7)
	case DateTypeMonth:
		prevDate = date.AddDate(0, -1, 0)
	case DateTypeQuarter:
		prevDate = date.AddDate(0, -3, 0)
	case DateTypeYear:
		prevDate = date.AddDate(-1, 0, 0)
	default:
		prevDate = date.AddDate(0, 0, -1)
	}

	prevStartTime, err := GetStartTime(dateType, prevDate)
	if err != nil {
		fmt.Printf("获取上一周期开始时间失败: %v\n", err)
		return
	}

	prevEndTime, err := GetEndTime(dateType, prevDate)
	if err != nil {
		fmt.Printf("获取上一周期结束时间失败: %v\n", err)
		return
	}

	// 格式化并打印时间
	fmt.Printf("上一%s的开始时间: %s\n", typeName, formatTime(prevStartTime, ""))
	fmt.Printf("上一%s的结束时间: %s\n\n", typeName, formatTime(prevEndTime, ""))
}

func showDaysBetween() {
	testDate := time.Now()
	fmt.Printf("测试日期: %s\n\n", formatTime(testDate, ""))
	// 获取当前周期的开始和结束时间
	startTime, err := GetStartTime(DateTypeMonth, testDate)
	if err != nil {
		fmt.Printf("获取开始时间失败: %v\n", err)
		return
	}

	endTime, err := GetEndTime(DateTypeMonth, testDate)
	if err != nil {
		fmt.Printf("获取结束时间失败: %v\n", err)
		return
	}

	result := GetDaysRange(startTime, endTime)
	fmt.Printf("日期范围：\n%s\n", strings.Join(result, "\n"))
}
