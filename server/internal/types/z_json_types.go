package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// LocalStringArray 字符串 数组类型
type LocalStringArray []string

// MarshalJSON 将 StringArray 编码为 JSON 格式
func (s LocalStringArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(s))
}

// UnmarshalJSON 将 JSON 解码为 StringArray
func (s *LocalStringArray) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*[]string)(s))
}

// Value 将 StringArray 转换为数据库驱动值
func (s LocalStringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan 从数据库读取数据并解码为 StringArray
func (s *LocalStringArray) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, s)
	}
	if str, ok := src.(string); ok {
		return json.Unmarshal([]byte(str), s)
	}
	return fmt.Errorf("cannot convert %v to StringArray", src)
}

// LocalIntArray int 数组类型
type LocalIntArray []uint64

// MarshalJSON 将 IntArray 编码为 JSON 格式
func (a LocalIntArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]uint64(a))
}

// UnmarshalJSON 将 JSON 解码为 IntArray
func (a *LocalIntArray) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*[]uint64)(a))
}

// Value 将 IntArray 转换为数据库驱动值（存储为 JSON）
func (a LocalIntArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan 从数据库读取数据并解码为 IntArray
func (a *LocalIntArray) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, a)
	}
	if str, ok := src.(string); ok {
		return json.Unmarshal([]byte(str), a)
	}
	return fmt.Errorf("cannot convert %v to IntArray", src)
}

// LocalJSONMap 定义自定义字典类型
type LocalJSONMap map[string]interface{}

// MarshalJSON 将 JSONMap 编码为 JSON 格式
func (m LocalJSONMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(m))
}

// UnmarshalJSON 将 JSON 解码为 JSONMap
func (m *LocalJSONMap) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*map[string]interface{})(m))
}

// Value 将 JSONMap 转换为数据库驱动值（存储为 JSON）
func (m LocalJSONMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan 从数据库读取数据并解码为 JSONMap
func (m *LocalJSONMap) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, m)
	}
	if str, ok := src.(string); ok {
		return json.Unmarshal([]byte(str), m)
	}
	return fmt.Errorf("cannot convert %v to JSONMap", src)
}
