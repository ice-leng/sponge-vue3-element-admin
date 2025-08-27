package cache

import (
	"testing"

	"golang.org/x/net/context"
)

func TestNewEnumCache(t *testing.T) {
	c := NewEnumCache()
	t.Log(c.GetAll(context.Background()))

	baseStatus, _ := c.Get(context.Background(), "base_status")
	t.Log(baseStatus)

	t.Log(c.GetLabel(context.Background(), "role_code", "ADMIN"))
}
