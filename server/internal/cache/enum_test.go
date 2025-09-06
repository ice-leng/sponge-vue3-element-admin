package cache

import (
	"admin/internal/config"
	"testing"

	"golang.org/x/net/context"
)

func TestNewEnumCache(t *testing.T) {
	_ = config.Init("")
	c := NewEnumCache()
	t.Log(c.GetAll(context.Background()))

	baseStatus, _ := c.Get(context.Background(), "base_status")
	t.Log(baseStatus)

	t.Log(c.GetLabel(context.Background(), "role_code", "ADMIN"))
}
