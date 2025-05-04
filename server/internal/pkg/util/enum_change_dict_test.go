package util

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func TestEnumChangeDict(t *testing.T) {
	// 获取当前文件的路径
	_, filename, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(path.Dir(filename)))
	enumDir := filepath.Join(root, "constant", "enum")
	result := EnumChangeDict(enumDir)
	t.Log(result)
}
