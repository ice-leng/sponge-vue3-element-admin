package datex

import (
	"errors"
	"time"
)

// 类型1日2周3月4季度5年
const (
	DateTypeDay     = iota + 1 // 日
	DateTypeWeek               // 周
	DateTypeMonth              // 月
	DateTypeQuarter            // 季度
	DateTypeYear               // 年
)

// CheckDateType 检查日期类型是否有效
// dateType: 类型参数（1=日，2=周，3=月，4=季度，5=年）
// 返回错误，如果类型无效
func checkDateType(dateType int) error {
	// 检查类型参数是否有效
	if dateType < DateTypeDay || dateType > DateTypeYear {
		return errors.New("无效的日期类型参数")
	}
	return nil
}

// GetDate 获取日期，如果为nil则返回当前时间
// date: 可选日期参数
// 返回处理后的时间
func getDate(date ...time.Time) time.Time {
	// 如果未提供日期，则使用当前系统时间
	now := time.Now()
	if len(date) > 0 {
		now = date[0]
	}
	return now
}

// GetStartTime 获取指定类型时间周期的开始时间
// dateType: 类型参数（1=日，2=周，3=月，4=季度，5=年）
// date: 可选日期参数，若为nil则使用当前系统时间
// 返回开始时间和可能的错误
func GetStartTime(dateType int, date ...time.Time) (time.Time, error) {
	// 检查类型参数是否有效
	if err := checkDateType(dateType); err != nil {
		return time.Time{}, err
	}

	// 获取日期
	currentDate := getDate(date...)

	// 获取当前时区的年、月、日信息
	year, month, day := currentDate.Date()
	hour, m, sec := 0, 0, 0 // 开始时间通常为0点0分0秒
	loc := currentDate.Location()
	switch dateType {
	case DateTypeDay: // 日
		// 当天的开始时间：当天的0点0分0秒
		return time.Date(year, month, day, hour, m, sec, 0, loc), nil

	case DateTypeWeek: // 周
		// 计算本周的开始时间（周一）
		weekday := currentDate.Weekday()
		if weekday == time.Sunday { // Go中Sunday是0，我们将其视为一周的最后一天
			weekday = 7
		}
		// 计算到周一的偏移天数
		offset := int(weekday) - 1
		// 本周一的日期
		weekStart := time.Date(year, month, day-offset, hour, m, sec, 0, loc)
		return weekStart, nil

	case DateTypeMonth: // 月
		// 本月的开始时间：本月1日的0点0分0秒
		return time.Date(year, month, 1, hour, m, sec, 0, loc), nil

	case DateTypeQuarter: // 季度
		// 计算当前季度的第一个月
		quarterMonth := ((int(month)-1)/3)*3 + 1
		// 本季度的开始时间：季度第一个月1日的0点0分0秒
		return time.Date(year, time.Month(quarterMonth), 1, hour, m, sec, 0, loc), nil

	case DateTypeYear: // 年
		// 本年的开始时间：1月1日的0点0分0秒
		return time.Date(year, 1, 1, hour, m, sec, 0, loc), nil

	default:
		// 默认返回当天的开始时间
		return time.Date(year, month, day, hour, m, sec, 0, loc), nil
	}
}

// GetEndTime 获取指定类型时间周期的结束时间
// date: 可选日期参数，若为nil则使用当前系统时间
// dateType: 类型参数（1=日，2=周，3=月，4=季度，5=年）
// 返回结束时间和可能的错误
func GetEndTime(dateType int, date ...time.Time) (time.Time, error) {
	// 检查类型参数是否有效
	if err := checkDateType(dateType); err != nil {
		return time.Time{}, err
	}

	// 获取日期
	currentDate := getDate(date...)

	// 获取当前时区的年、月、日信息
	year, month, day := currentDate.Date()
	hour, m, sec := 23, 59, 59 // 结束时间通常为23点59分59秒
	loc := currentDate.Location()

	switch dateType {
	case DateTypeDay: // 日
		// 当天的结束时间：当天的23点59分59秒
		return time.Date(year, month, day, hour, m, sec, 999999999, loc), nil

	case DateTypeWeek: // 周
		// 计算本周的开始时间（周一）
		weekday := currentDate.Weekday()
		if weekday == time.Sunday { // Go中Sunday是0，我们将其视为一周的最后一天
			weekday = 7
		}
		// 计算到周一的偏移天数
		offset := int(weekday) - 1
		// 本周日的日期（周一+6天）
		weekEnd := time.Date(year, month, day-offset+6, hour, m, sec, 999999999, loc)
		return weekEnd, nil

	case DateTypeMonth: // 月
		// 计算下个月的第一天
		nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, loc)
		// 本月的最后一天是下个月第一天的前一天
		lastDay := nextMonth.AddDate(0, 0, -1)
		// 本月的结束时间：本月最后一天的23点59分59秒
		return time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), hour, m, sec, 999999999, loc), nil

	case DateTypeQuarter: // 季度
		// 计算当前季度的第一个月
		quarterMonth := ((int(month)-1)/3)*3 + 1
		// 计算下一个季度的第一天
		nextQuarter := time.Date(year, time.Month(quarterMonth+3), 1, 0, 0, 0, 0, loc)
		// 本季度的最后一天是下一个季度第一天的前一天
		lastDay := nextQuarter.AddDate(0, 0, -1)
		// 本季度的结束时间：本季度最后一天的23点59分59秒
		return time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), hour, m, sec, 999999999, loc), nil

	case DateTypeYear: // 年
		// 本年的结束时间：12月31日的23点59分59秒
		return time.Date(year, 12, 31, hour, m, sec, 999999999, loc), nil

	default:
		// 默认返回当天的结束时间
		return time.Date(year, month, day, hour, m, sec, 999999999, loc), nil
	}
}

// GetLastStartTime 获取上一个时间周期的开始时间
// date: 可选日期参数，若为nil则使用当前系统时间
// dateType: 类型参数（1=日，2=周，3=月，4=季度，5=年）
// 返回上一个时间周期的开始时间和可能的错误
func GetLastStartTime(dateType int, date time.Time) (time.Time, error) {
	// 检查类型参数是否有效
	if err := checkDateType(dateType); err != nil {
		return time.Time{}, err
	}

	// 获取日期
	currentDate := getDate(date)

	// 根据类型计算上一个周期的日期
	var prevDate time.Time

	switch dateType {
	case DateTypeDay: // 日
		// 前一天
		prevDate = currentDate.AddDate(0, 0, -1)

	case DateTypeWeek: // 周
		// 前一周
		prevDate = currentDate.AddDate(0, 0, -7)

	case DateTypeMonth: // 月
		// 前一个月
		prevDate = currentDate.AddDate(0, -1, 0)

	case DateTypeQuarter: // 季度
		// 前一个季度
		prevDate = currentDate.AddDate(0, -3, 0)

	case DateTypeYear: // 年
		// 前一年
		prevDate = currentDate.AddDate(-1, 0, 0)

	default:
		// 默认前一天
		prevDate = currentDate.AddDate(0, 0, -1)
	}

	// 获取上一个周期的开始时间
	return GetStartTime(dateType, prevDate)
}

// GetLastEndTime 获取上一个时间周期的结束时间
// date: 可选日期参数，若为nil则使用当前系统时间
// dateType: 类型参数（1=日，2=周，3=月，4=季度，5=年）
// 返回上一个时间周期的结束时间和可能的错误
func GetLastEndTime(dateType int, date time.Time) (time.Time, error) {
	// 检查类型参数是否有效
	if err := checkDateType(dateType); err != nil {
		return time.Time{}, err
	}

	// 获取日期
	currentDate := getDate(date)

	// 根据类型计算上一个周期的日期
	var prevDate time.Time

	switch dateType {
	case DateTypeDay: // 日
		// 前一天
		prevDate = currentDate.AddDate(0, 0, -1)

	case DateTypeWeek: // 周
		// 前一周
		prevDate = currentDate.AddDate(0, 0, -7)

	case DateTypeMonth: // 月
		// 前一个月
		prevDate = currentDate.AddDate(0, -1, 0)

	case DateTypeQuarter: // 季度
		// 前一个季度
		prevDate = currentDate.AddDate(0, -3, 0)

	case DateTypeYear: // 年
		// 前一年
		prevDate = currentDate.AddDate(-1, 0, 0)

	default:
		// 默认前一天
		prevDate = currentDate.AddDate(0, 0, -1)
	}

	// 获取上一个周期的结束时间
	return GetEndTime(dateType, prevDate)
}

// IsToday 判断指定日期是否为今天
// dateType: 类型参数（1=日，2=周，3=月，4=季度，5=年）
// date: 要判断的日期，若为nil则使用当前系统时间
// 返回布尔值表示是否为今天，以及可能的错误
func IsToday(dateType int, date time.Time) (bool, error) {
	// 检查日期类型是否有效
	if err := checkDateType(dateType); err != nil {
		return false, err
	}

	// 获取日期
	currentDate := getDate(date)

	// 获取当前日期类型的开始时间
	dateTypeDate, err := GetStartTime(dateType)
	if err != nil {
		return false, err
	}

	// 获取当前日期的年、月、日
	year, month, day := dateTypeDate.Date()

	// 获取指定日期的年、月、日
	dYear, dMonth, dDay := currentDate.Date()

	// 比较年、月、日是否相同
	return year == dYear && month == dMonth && day == dDay, nil
}

// GetDaysRange
// 给定2个日期范围 返回 这个2个日期之间的所有日
func GetDaysRange(startDate, endDate time.Time, format ...string) []string {
	// 确保开始日期不晚于结束日期
	if startDate.After(endDate) {
		startDate, endDate = endDate, startDate
	}

	// 将时间调整为当天的0点0分0秒
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	// 计算日期差
	duration := endDate.Sub(startDate)
	days := int(duration.Hours()/24) + 1 // 包括开始和结束日期

	// 创建结果切片
	result := make([]string, 0, days)
	f := "2006-01-02"
	if len(format) > 0 {
		f = format[0]
	}

	// 逐天添加日期
	currentDate := startDate
	for i := 0; i < days; i++ {
		result = append(result, currentDate.Format(f))
		currentDate = currentDate.AddDate(0, 0, 1) // 增加一天
	}

	return result
}
