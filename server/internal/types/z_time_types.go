package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type LocalDateTime time.Time

// MarshalJSON 将 LocalDateTime 序列化为 JSON 时间字符串
// 使用本地时区格式化时间
func (dt LocalDateTime) MarshalJSON() ([]byte, error) {
	value := ""
	if !time.Time(dt).IsZero() {
		// 确保使用本地时区格式化时间
		loc, err := time.LoadLocation("Local")
		if err != nil {
			// 如果无法加载本地时区，使用 UTC
			loc = time.UTC
		}
		value = time.Time(dt).In(loc).Format("2006-01-02 15:04:05")
	}
	return []byte(fmt.Sprintf(`"%s"`, value)), nil
}

// UnmarshalJSON 解析 JSON 时间字符串为 LocalDateTime
// 使用本地时区解析时间，避免时区偏移问题
func (dt *LocalDateTime) UnmarshalJSON(b []byte) error {
	// 检查是否为 JSON null
	if string(b) == "null" {
		return nil
	}

	// 解析带引号的 JSON 时间字符串
	timeStr := strings.Trim(string(b), "\"")

	// 使用本地时区解析时间
	loc, err := time.LoadLocation("Local")
	if err != nil {
		// 如果无法加载本地时区，使用 UTC
		loc = time.UTC
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
	if err != nil {
		return err
	}

	*dt = LocalDateTime(t) // 直接赋值给 *dt
	return nil
}

// Value 将 LocalDateTime 转换为数据库驱动值（存储为 JSON）
func (dt *LocalDateTime) Value() (driver.Value, error) {
	return dt.MarshalJSON()
}

// Scan 从数据库读取数据并解码为 LocalDateTime
func (dt *LocalDateTime) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*dt = LocalDateTime(t)
		return nil
	}
	if bytes, ok := src.([]byte); ok {
		return dt.UnmarshalJSON(bytes)
	}
	if str, ok := src.(string); ok {
		return dt.UnmarshalJSON([]byte(str))
	}
	return fmt.Errorf("cannot convert %v to LocalDateTime", src)
}

type LocalDate time.Time

func (dt LocalDate) MarshalJSON() ([]byte, error) {
	value := ""
	if !time.Time(dt).IsZero() {
		// 确保使用本地时区格式化时间
		loc, err := time.LoadLocation("Local")
		if err != nil {
			// 如果无法加载本地时区，使用 UTC
			loc = time.UTC
		}
		value = time.Time(dt).In(loc).Format("2006-01-02")
	}
	return []byte(fmt.Sprintf(`"%s"`, value)), nil
}

func (dt *LocalDate) UnmarshalJSON(b []byte) error {
	// 检查是否为 JSON null
	if string(b) == "null" {
		return nil
	}

	// 解析带引号的 JSON 时间字符串
	timeStr := strings.Trim(string(b), "\"")

	// 使用本地时区解析时间
	loc, err := time.LoadLocation("Local")
	if err != nil {
		// 如果无法加载本地时区，使用 UTC
		loc = time.UTC
	}

	// 解析带引号的 JSON 时间字符串
	t, err := time.ParseInLocation(`2006-01-02`, timeStr, loc)
	if err != nil {
		return err
	}

	*dt = LocalDate(t) // 直接赋值给 *dt
	return nil
}

// Value 将 LocalDate 转换为数据库驱动值（存储为 JSON）
func (dt *LocalDate) Value() (driver.Value, error) {
	return dt.MarshalJSON()
}

// Scan 从数据库读取数据并解码为 LocalDate
func (dt *LocalDate) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		*dt = LocalDate(t)
		return nil
	}
	if bytes, ok := src.([]byte); ok {
		return dt.UnmarshalJSON(bytes)
	}
	if str, ok := src.(string); ok {
		return dt.UnmarshalJSON([]byte(str))
	}
	return fmt.Errorf("cannot convert %v to LocalDate", src)
}
