package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

var (
	// 当前代理服务器的启动端口
	port = "8083"
	// 被代理的服务器地址
	proxy_addr = "http://127.0.0.1:8081"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// 解析上游的url
	URL, err := url.Parse(proxy_addr)
	if err != nil {
		fmt.Printf("Fail to parse proxy_addr error:[%v]\n", err)
		return
	}
	// 偷天换日
	// 原来这里踩坑：[http: no Host in request URL] 注意点：我们得将Host赋值给r.URL.Host
	r.URL.Scheme = URL.Scheme
	r.URL.Host = URL.Host

	// 将请求转发到下游
	transport := http.DefaultTransport
	res, err := transport.RoundTrip(r)
	if err != nil {
		fmt.Printf("Fail to get response from transport roundTrip error:[%v]\n", err)
		return
	}

	// 将下游的响应中的内容赋值给上游的responseWriter
	for k, value := range res.Header {
		for _, v := range value {
			w.Header().Add(k, v)
		}
	}
	// 关闭下游返回的响应体
	defer res.Body.Close()
	bufio.NewReader(res.Body).WriteTo(w)
}

// 反向代理服务器
func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))
}
