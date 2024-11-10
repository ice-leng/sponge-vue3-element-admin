package util

import (
	"unicode"
)

func UcFirst(str string) string {
	runes := []rune(str)
	if len(runes) == 0 {
		return str
	}
	// 将第一个字符转换为大写
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
