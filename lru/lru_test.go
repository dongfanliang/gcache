package lru

import (
	"testing"
)

func TestGet(t *testing.T) {
	lru := New(0)

	lru.Add("1", "111")
	lru.Add("2", "222")
	lru.Add("3", "333")

	v1, ok := lru.Get("1")
	t.Log(v1, ok)

	v3, ok := lru.Get("3")
	t.Log(v3, ok)
	t.Log(lru.Len(), lru.UsedBytes())
}

func TestRemove(t *testing.T) {
	lru := New(12)

	lru.Add("1", "111")
	lru.Add("2", "222")
	lru.Add("3", "333")

	v1, ok := lru.Get("1")
	t.Log(v1, ok)

	lru.Add("4", "444")
	v1, ok = lru.Get("2")
	t.Log(v1, ok)
	t.Log(lru.Len(), lru.UsedBytes(), lru.Keys())

	lru.Remove("1")
	t.Log(lru.Len(), lru.UsedBytes(), lru.Keys())
}

func TestEvict(t *testing.T) {
	evictedKeys := make([]string, 0)
	onEvictedFun := func(key string, value string) {
		evictedKeys = append(evictedKeys, key)
	}
	lru := New(12)
	lru.OnEvicted = onEvictedFun

	lru.Add("1", "111")
	lru.Add("2", "222")
	lru.Add("3", "333")
	lru.Add("4", "444")
	lru.Add("5", "555")

	t.Log(evictedKeys)
	t.Log(lru.Len(), lru.UsedBytes(), lru.Keys())
}
