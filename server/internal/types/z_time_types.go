package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type LocalDateTime time.Time

func (dt LocalDateTime) MarshalJSON() ([]byte, error) {
	value := ""
	if !time.Time(dt).IsZero() {
		value = time.Time(dt).Format("2006-01-02 15:04:05")
	}
	return []byte(fmt.Sprintf(`"%s"`, value)), nil
}

func (dt *LocalDateTime) UnmarshalJSON(b []byte) error {
	// 检查是否为 JSON null
	if string(b) == "null" {
		return nil
	}

	// 解析带引号的 JSON 时间字符串
	t, err := time.Parse(`2006-01-02 15:04:05`, strings.Trim(string(b), "\""))
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
		value = time.Time(dt).Format("2006-01-02 15:04:05")
	}
	return []byte(fmt.Sprintf(`"%s"`, value)), nil
}

func (dt *LocalDate) UnmarshalJSON(b []byte) error {
	// 检查是否为 JSON null
	if string(b) == "null" {
		return nil
	}

	// 解析带引号的 JSON 时间字符串
	t, err := time.Parse(`2006-01-02 15:04:05`, strings.Trim(string(b), "\""))
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
	if bytes, ok := src.([]byte); ok {
		return dt.UnmarshalJSON(bytes)
	}
	if str, ok := src.(string); ok {
		return dt.UnmarshalJSON([]byte(str))
	}
	return fmt.Errorf("cannot convert %v to LocalDate", src)
}
