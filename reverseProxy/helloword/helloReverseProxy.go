package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	proxyServerAddr = ":8085"
	targetAddr = "http://127.0.0.1:8081/base"
)

// 最终实现的效果：
/**
 *	用户浏览器访问：http://127.0.0.1:8085/qwer
 *	请求被转发到：http://127.0.0.1:8081/base/qwer
 */

func main() {
	// 解析targetAddr
	URL, err := url.Parse(targetAddr)
	if err != nil {
		fmt.Printf("Fail to parse url:[%v] error:[%v]", targetAddr,err)
		return
	}
	// 构建reverseProxy
	proxy := httputil.NewSingleHostReverseProxy(URL)
	// 启动服务
	log.Fatal(http.ListenAndServe(proxyServerAddr,proxy))
}

