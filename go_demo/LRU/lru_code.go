package LRU

import "container/list"

//创建一个cache数据结构，使用go语言标准库中的list.List创建双向链表

type Cache struct {
	maxBytes int64		//代表了最大使用内存
	nbytes   int64		//代表了当前已经使用的内存
	ll       *list.List		//建立的双向链表
	cache    map[string]*list.Element		//创建的字典，分别指向了对应的节点的位置
	OnEvicted func(key string, value Value)		//回调函数，不知道干啥的
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func (c *Cache) Len() int{
	return c.ll.Len()
}

//Cache实例化
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//查找功能，从字典中找到对应的双向链表节点
func (c *Cache) Get(key string) (value Value, ok bool){
	element , ok :=c.cache[key]
	if  ok {		//存在
		c.ll.MoveToFront(element)		//放到最前面，即放到队尾
		kv := element.Value.(*entry)	//获取对应的元素信息
		return kv.value, true
	}
	return nil, false
}

//删除功能，把最近最少访问的节点淘汰掉，即淘汰掉队首元素
func (c *Cache) RemoveOldest(){
	element := c.ll.Back()		//取到队首节点，从链表中删除
	if element != nil{
		c.ll.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)		//删除尾节点
		c.nbytes -=int64(len(kv.key)) + int64(kv.value.Len())	//把占用的内存删除
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key,kv.value)
		}
	}
}

//新增，修改
func (c *Cache) Add (key string, value Value){
	element,ok :=c.cache[key]
	if ok {
		c.ll.MoveToFront(element)
		kv := element.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		element := c.ll.PushFront(&entry{key, value})
		c.cache[key] = element
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes !=0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}

}

