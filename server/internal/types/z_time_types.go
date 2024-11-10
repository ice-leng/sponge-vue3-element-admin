package types

import (
	"fmt"
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
	t, err := time.Parse(`"2006-01-02 15:04:05"`, string(b)) // 包含引号
	if err != nil {
		return err
	}

	*dt = LocalDateTime(t) // 直接赋值给 *dt
	return nil
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
	t, err := time.Parse(`"2006-01-02 15:04:05"`, string(b)) // 包含引号
	if err != nil {
		return err
	}

	*dt = LocalDate(t) // 直接赋值给 *dt
	return nil
}
