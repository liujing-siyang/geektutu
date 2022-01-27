package lru

import "container/list"

type Cache struct {
	maxBytes  int64                    //允许使用的最大内存
	nbytes    int64                    //当前已使用的内存
	cacheList *list.List               //双向链表，存放所有缓存值
	cacheMap  map[string]*list.Element //键是字符串，值是双向链表中对应节点的指针

	OnEvicted func(key string, value Value) //回调函数，删除操作时的定制需求
}

//双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		cacheList: list.New(),
		cacheMap:  make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cacheMap[key]; ok {
		c.cacheList.MoveToFront(ele) //缓存链表中存在，则将其放到队头位置
		kv := ele.Value.(*entry)     //Element的Value为空接口，所以断言成我们存入的节点类型entry
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.cacheList.Back() //取队尾元素
	if ele != nil {
		c.cacheList.Remove(ele) //链表中删除
		kv := ele.Value.(*entry)
		delete(c.cacheMap, kv.key)                             //map中删除
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) //更新当前使用的内存
		if c.OnEvicted != nil {                                //如果回调函数 OnEvicted 不为 nil，则调用回调函数
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	//缓存链表中存在,修改
	if ele, ok := c.cacheMap[key]; ok {
		c.cacheList.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { //不存在，则添加
		ele := c.cacheList.PushFront(&entry{key, value})//往c.cacheList添加的为指针，所以Get获取时强转成也是指针
		c.cacheMap[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.nbytes > c.maxBytes { //c.maxBytes == 0,就没有缓存能力了，那当Add一个元素后马上就要删除，c.maxBytes != 0这个条件不用加
		c.RemoveOldest()
	}
}

func (c Cache) Len() int {
	return c.cacheList.Len()
}
