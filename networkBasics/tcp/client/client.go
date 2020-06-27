package main

import (
	"fmt"
	"gateway-plus/networkBasics/tcp/coder"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":9090")
	defer conn.Close()
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}

	// 将数据编码并发送出去
	coder.Encode(conn,"hi server i am here");

	time.Sleep(time.Second*10)
}
