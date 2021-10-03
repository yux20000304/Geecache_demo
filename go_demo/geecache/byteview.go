package geecache

//建立一个数据结构表示缓存值，注意这里只实现了读操作
type ByteView struct {
	b [] byte
}

//返回长度
func (v ByteView) Len() int{
	return len(v.b)
}

//返回v的切片
func (v ByteView) ByteSlice() []byte{
	return cloneBytes(v.b)
}

//返回v的string类型
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}