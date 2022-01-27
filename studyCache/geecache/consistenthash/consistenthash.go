package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           //hash函数
	replicas int            //虚拟节点倍数,通过虚拟节点解决数据倾斜的问题
	keys     []int          //节点编号组成的哈希环
	hashMap  map[int]string //节点编号对应缓存服务器的地址URL
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m * Map)Add(urls ...string){
	for _,url := range urls{
		for i := 0;i< m.replicas;i++{
			hash := int(m.hash([]byte(strconv.Itoa(i) + url)))//存在hash碰撞的情况，如果两个不同的url处理后hash值一样，m.hashMap中对应关系就会被覆盖
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = url
		}
	}
	sort.Ints(m.keys)
}

func (m * Map)Get(key string)string{//key是要检索的数据
	if len(m.keys) == 0{
		return ""
	}
	hash := int(m.hash([]byte(key)))
	//在[0,n)范围内，返回使得f(i)为true的的索引i,如果没有这样的则返回n,所以index的值在[0,n]中
	idx := sort.Search(len(m.keys),func(i int) bool {
		return m.keys[i] >= hash
	})
	//idx == len(m.keys)，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
	return m.hashMap[m.keys[idx%len(m.keys)]]
}