package Consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 创建一个依赖注入函数
type Hash func(data []byte) uint32

// 定义一个map数据结构
type Map struct {
	hash     Hash		//使用的hash函数
	replicas int		//虚拟节点的倍数
	keys     []int 		//hash环的keys
	hashMap  map[int]string		//虚拟节点和真实节点的hash表
}

// 创建一个实例对象
func New(replicas int, fn Hash) *Map {
	m := &Map{					//创建实例对象
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE		//默认使用crcchecksum这个hash函数
	}
	return m
}

// 添加节点的方法，可以传入多个节点的名称
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// 顺时针寻找到最近的节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	// Binary search for appropriate replica.
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}