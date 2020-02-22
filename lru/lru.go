package lru

import "container/list"

type entry struct {
	key   string
	value string
}

type Cache struct {
	maxBytes  int64 // 最大使用内存，0为不限制
	nbytes    int64 // 已使用内存
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value string) // 删除  回调函数
}

func New(maxBytes int64) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

func (c *Cache) Get(key string) (value string, ok bool) {
	if c.cache == nil {
		return
	}

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) Add(key string, value string) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		kv.value = value
		return
	}

	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	c.nbytes += int64(len(key)) + int64(len(value))

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		// 添加的时候才会触发
		c.removeOldest()
	}
}

func (c *Cache) Remove(key string) {
	if c.cache == nil {
		return
	}
	if e, ok := c.cache[key]; ok {
		c.removeElement(e)
	}
}

func (c *Cache) removeOldest() {
	if c.cache == nil {
		return
	}

	e := c.ll.Back()
	if e != nil {
		c.removeElement(e)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	c.nbytes -= int64(len(kv.key)) + int64(len(kv.value))
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

func (c *Cache) Keys() (keys []string) {
	if c.cache == nil {
		return
	}

	keys = make([]string, 0, c.Len())
	for k, _ := range c.cache {
		keys = append(keys, k)
	}
	return keys
}

func (c *Cache) UsedBytes() int64 {
	if c.cache == nil {
		return 0
	}
	return c.nbytes
}

func (c *Cache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.ll = nil
	c.cache = nil
}
