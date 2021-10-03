package geecache

import (
	"go_demo/LRU"
	"sync"
)
//完整的一个cache
type cache struct{
	mutex1 sync.Mutex
	lru *LRU.Cache
	cacheBytes int64
}

func (c *cache) add (key string, value ByteView)  {
	c.mutex1.Lock()
	defer c.mutex1.Unlock()
	if c.lru == nil {
		c.lru = LRU.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get (key string) (value ByteView, ok bool){
	c.mutex1.Lock()
	defer c.mutex1.Unlock()
	if c.lru == nil{
		return
	}
	if v,ok := c.lru.Get(key); ok{
		return v.(ByteView), ok
	}
	return
}

