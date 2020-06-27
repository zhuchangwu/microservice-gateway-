package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type myProxy struct {
}



func (p *myProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Recieve request, requesrMethod:[%v]  requestHost:[%v]  resquestRemoteAddr:[%v]\n", r.Method, r.Host, r.RemoteAddr)
	// 获取http.DefaultTranport实例，原来有写文章具体的看过这个Tranport。
	// 他通过roundTrip，可以获取和server的conn，并将我们传递给他的req发送到server，并获取到响应

	tranport := http.DefaultTransport

	// 为了防止我们的给req添加新属性时对req造成影响，故前拷贝一份req，再添加新属性
	// todo 这是深度拷贝～
	req := new(http.Request)
	*req = *r

	if host, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := req.Header["X-Forward-For"]; ok {
			host = strings.Join(prior, ",") + "," + host
		}
		// 更新头信息
		req.Header.Set("X-Forward-For", host)
	}

	// 将请求转发到下游
	res, err := tranport.RoundTrip(req)
	if err != nil {
		fmt.Printf("Fail to tranport roundTrip err:[%v]", err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	// 将下游返回过来的头信息添加到上游请求对应的header中
	for resHeaderKey, resHeaderValue := range res.Header {
		for _, v := range resHeaderValue {
			w.Header().Add(resHeaderKey, v)
		}
	}
	// 将请求写会到上游
	w.WriteHeader(res.StatusCode)
	io.Copy(w,res.Body)
	res.Body.Close()
}

// 正向代理服务器
func main() {
	// 创建http服务器，并这个server注册路由
	http.Handle("/", &myProxy{})
	http.ListenAndServe(":8080", nil)
}
