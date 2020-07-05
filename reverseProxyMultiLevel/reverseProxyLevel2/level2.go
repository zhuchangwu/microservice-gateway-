package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	// 第二重代理服务器的地址
	SERVER_ADDR = "127.0.0.1:8086"

	// 第二重代理服务器的下游真实服务器地址
	// 通过一个随机的负载均衡策略，将请求转发到真正的服务提供者上
	PROXY_ADDR1 = "http://127.0.0.1:8081/base"
	PROXY_ADDR2 = "http://127.0.0.1:8082/base"

)

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       10 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

/**
 *	基于ReverseProxy原生代码改写
 *	如果：请求路径是 http://localhost:8080/dir
 *	如果：代理路径是 http://localhost:8080/base
 *	最终的路径是 http://localhost:8080/base/dir
 */
func NewMultiHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {

	// 调度者
	director := func(req *http.Request) {
		reg, err := regexp.Compile("^dir(.*)")
		if err != nil {
			fmt.Printf("Fail to complie reg err:[%v]", err)
			return
		}
		req.URL.Path = reg.ReplaceAllString(req.URL.Path, "$1")

		// todo 随机的负载均衡
		targetIndex := rand.Intn(len(targets))
		target := targets[targetIndex]
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		// url地址重写
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		// todo 只在第一重代理中设置它的header头
		req.Header.Set("X-Real-Ip", req.RemoteAddr)
		// todo X-Forward-For这个请求头在ReverseProxy已经为我们设置好了 （reverseproxy.go的ServeHTTP方法中已经为我们设置好了）

	}
	// 声明ModifyResponse类型的匿名函数
	modifyFunc := func(res *http.Response) error {
		if res.StatusCode != 200 {

			// 获取到下游返回到res
			oldPayLoad, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("Fail to read response body error:[%v]", err)
				return err
			}

			// 在获取到响应前追加一部分内容
			newRes := []byte("hello " + string(oldPayLoad))

			//将更改后的body重新写会responseBody中
			/**
			res.Body类型为ReadCloser如下，想为Body赋值看起腰要去实现下面的接口，重写里面的方法
			但是其实是不用这么搞的～，有现成的工具类
			type ReadCloser interface {
			Reader
			Closer
			*/
			res.Body = ioutil.NopCloser(bytes.NewBuffer(newRes))
			// 复写content-length
			res.ContentLength = int64(len(newRes))
			// 添加响应头，告诉客户端content-length
			res.Header.Set("Content-Length", fmt.Sprint(len(newRes)))
		}
		return nil
	}

	errorHandler := func(w http.ResponseWriter,e *http.Request,err error) {
		fmt.Printf("Error happened err:[%v]",err)
	}

	return &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyFunc,
		Transport:      transport,
		ErrorHandler:   errorHandler,
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}


func main() {
	// 实际的服务提供者
	URL1, err := url.Parse(PROXY_ADDR1)
	if err != nil {
		fmt.Printf("Fail to parse PROXY_ADDR:[%v] error:[%v}", PROXY_ADDR1, err)
		return
	}
	// 实际的服务提供者
	URL2, err := url.Parse(PROXY_ADDR2)
	if err != nil {
		fmt.Printf("Fail to parse PROXY_ADDR:[%v] error:[%v}", PROXY_ADDR2, err)
		return
	}
	urls:=[]*url.URL{URL1,URL2}
	proxy := NewMultiHostsReverseProxy(urls)
	log.Fatal(http.ListenAndServe(SERVER_ADDR, proxy))
}


