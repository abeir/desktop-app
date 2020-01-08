package core

import (
	"bytes"
)

// MakeStringStack 创建 StringStarck
func MakeStringStack() *StringStack{
	return &StringStack{}
}

// StringStarck 是一个模拟的字符串栈，提供了入栈出栈等操作
type StringStack struct {
	size int64
	data bytes.Buffer
}

// Push 将一个字符串入栈
// param
//    data: 准备入栈的字符串
// return
//    int: 成功入栈后返回入栈的字符串长度
//    error: 入栈失败返回错误
func (s *StringStack) Push(data string) (int, error) {
	n, err := s.data.WriteString(data)
	s.size += int64(n)
	return n, err
}

// PushByte 将一个字节入栈
// param
//    data: 准备入栈的字节
// return
//    int: 成功入栈后返回入栈的字节长度
//    error: 入栈失败返回错误
func (s *StringStack) PushByte(data byte) (int, error) {
	err := s.data.WriteByte(data)
	s.size += 1
	return 1, err
}

// PushBytes 将一个字节切片入栈
// param
//    data: 准备入栈的字节切片
// return
//    int: 成功入栈后返回入栈的字节切片长度
//    error: 入栈失败返回错误
func (s *StringStack) PushBytes(data []byte) (int, error) {
	n, err := s.data.Write(data)
	s.size += int64(len(data))
	return n, err
}

// Pop 将当前栈中的所有数据出栈，并清空当前栈
// return
//    string: 栈中存储的字符串
//    error: 出栈失败的错误
func (s *StringStack) Pop() (string, error){
	data := make([]byte, s.size)
	_, err := s.data.Read(data)
	s.Clear()
	return string(data), err
}

// Size 当前栈的容量
func (s *StringStack) Size() int64{
	return s.size
}

// Clear 清空当前栈
func (s *StringStack) Clear(){
	s.data.Reset()
	s.size = 0
}




