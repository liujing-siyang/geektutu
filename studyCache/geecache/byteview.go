package geecache

type ByteView struct{
	b []byte//b 是只读的，使用 ByteSlice() 方法返回一个拷贝，防止缓存值被外部程序修改
}

func (v ByteView)Len()int{
	return len(v.b)
}

func (v ByteView)ByteSlice()[]byte{
	return cloneBytes(v.b)
}

func (v ByteView)String() string{
	return string(v.b)
}

func cloneBytes(b []byte)[]byte{
	c := make([]byte,len(b))
	copy(c,b)
	return c
}