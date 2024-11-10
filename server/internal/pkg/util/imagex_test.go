package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidURL(t *testing.T) {
	url := "https://www.baidu.com"
	is := IsValidURL(url)
	assert.Equal(t, true, is)

	str := "www.baidu.com"
	is = IsValidURL(str)
	assert.Equal(t, false, is)
}

func TestImageMakePath(t *testing.T) {
	host := "https://www.baidu.com/"
	path := "/a.jpg"

	url := ImageMakePath(path, host)
	assert.Equal(t, "https://www.baidu.com/a.jpg", url)

	str := "a.jpg"
	url = ImageMakePath(str, host)
	assert.Equal(t, "https://www.baidu.com/a.jpg", url)
}
