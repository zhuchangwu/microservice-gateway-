package main

import "fmt"

// f2是一个普通函数，有两个入参数
func f2() {
	fmt.Printf("f2222")
}

// f1函数的入参是一个f2类型的函数
func f1(f2 func()) {
	f2()
}

func main() {
	f1(f2)
}
