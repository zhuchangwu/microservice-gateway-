package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// 创建路由器
	// 为路由器绑定路由规则
	// 创建服务器
	// 监听端口，启动读服务
	mux := http.NewServeMux()
	mux.HandleFunc("/login", doLogin)

	server := &http.Server{
		Addr:         ":8081",
		WriteTimeout: time.Second * 2,
		Handler:      mux,
	}
	log.Fatal(server.ListenAndServe())
}

func doLogin(writer http.ResponseWriter,req *http.Request){
	_, err := writer.Write([]byte("do login"))
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
}