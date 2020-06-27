package main

import "fmt"

// 自定义类型，handler本质上是一个函数
type HandlerFunc func(string, int)
type myInt int

// 这应该是一个闭bao
func (f HandlerFunc) Serve(name string, age int) {
	f(name, age)
}

// 具体的处理函数
func HelloHandle(name string, age int) {
	fmt.Printf("name:[%v] age:[%v]", name, age)
}


func main() {
	handlerFunc := HandlerFunc(HelloHandle)
	handlerFunc.Serve("tom", 12)

	m := myInt(12)
	fmt.Printf("\n myint:[%v]", m)
}
