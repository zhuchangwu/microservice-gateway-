package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {
	// 创建连接池
	// 创建客户端，绑定连接池
	// 发送请求
	// 读取响应
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   3000 * time.Second, // 连接超时
			KeepAlive: 3000 * time.Second, // 长连接存活的时间
		}).DialContext,
		MaxIdleConns:          100,              // 最大空闲连接数
		IdleConnTimeout:       100 * time.Second,
		TLSHandshakeTimeout:   100 * time.Second, // tls握手超时时间
		ExpectContinueTimeout: 100 * time.Second,  // 100-continue 状态码超时时间
	}

	// 创建客户端
	client := &http.Client{
		Timeout:   time.Second * 10, //请求超时时间
		Transport: transport,
	}

	// 请求数据
	res, err := client.Get("http://localhost:8081/login")
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
	fmt.Printf("Read from http server res:[%v]", string(bytes))
}
