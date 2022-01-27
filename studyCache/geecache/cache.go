package geecache

import (
	"studycache/geecache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex//是不是应该换成读写锁，add时使用写锁，get使用读锁
	//cache 的 get 和 add 都涉及到写操作(LRU 将最近访问元素移动到链表头)，所以不能直接改为读写锁。
	//如果 cache 侧和 LRU 侧同时使用锁细颗粒度控制，是有优化空间的，可以尝试下
	lru        *lru.Cache
	cacheBytes int64 //为了在延迟初始化lru.Cache时决定缓存的最大值
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//c.lru == nil,应该要报错，否则返回空值和false，不清楚是缓存中没找到还是根本还没有这个缓存
	//后续直接取查本地缓存了，有的话会加入cache,调用add完成了c.lru的初始化，这样的话就无需区分了
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return ByteView{}, false
}
