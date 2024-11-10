package util

import (
	"fmt"
	"net/url"
	"strings"
)

// IsValidURL  is url
func IsValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	// 检查是否存在解析错误以及是否包含协议（如 http 或 https）
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}

// ImageMakePath   host + path
func ImageMakePath(path, host string) string {
	if path == "" || IsValidURL(path) {
		return path
	}

	return fmt.Sprintf("%s/%s", strings.TrimRight(host, "\t\n\r\x00\x0B/"), strings.TrimLeft(path, "/"))
}
