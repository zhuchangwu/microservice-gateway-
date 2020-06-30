package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

/**
 *	基于ReverseProxy原生代码改写
 *	如果：请求路径是 http://localhost:8080/dir
 *	如果：代理路径是 http://localhost:8080/base
 *	最终的路径是 http://localhost:8080/base/dir
 */
func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	// http:/?name=123
	// RawQuery: name=123
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		// Scheme: http
		req.URL.Scheme = target.Scheme
		// Host: localhost:8080
		req.URL.Host = target.Host
		// singleJoiningSlash拼接 targetPath 和 req.URL.Path
		// 最终结果是：http:xxx/targetPath/req.URL.Path
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

	return &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyFunc,
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

var (
	PROXY_ADDR = "http://127.0.0.1:8081/base" // 代理的地址
	// todo 地址不要写成下面的样子～～～
	// SERVER_ADDR = "http://127.0.0.1:8085/dir" // 当前代理服务器的地址
	SERVER_ADDR = "127.0.0.1:8085" // 当前代理服务器的地址
)

func main() {
	// 思路：将代理地址
	URL, err := url.Parse(PROXY_ADDR)
	if err != nil {
		fmt.Printf("Fail to parse PROXY_ADDR:[%v] error:[%v}", PROXY_ADDR, err)
		return
	}
	// 这里使用reverProxy不是httputil中原生的，而是我们改写的，修改点：我们为proxy添加上了ModifyResponse
	proxy := NewSingleHostReverseProxy(URL)
	log.Fatal(http.ListenAndServe(SERVER_ADDR, proxy))
}

// 测试结果：
/**
	MacBook-Pro% curl http://localhost:8085/dir
	hello http://127.0.0.1:8081/base/dir
*/
