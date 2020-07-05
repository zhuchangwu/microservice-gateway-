package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RealServer struct {
	Addr string
}

func (r *RealServer) RUN() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.IndexHandler)
	mux.HandleFunc("/base/err", r.ErrorHandler)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: 2 * time.Second,
		Handler:      mux,
	}
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

func (rs *RealServer) IndexHandler(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("http://%s%s\n", rs.Addr, r.URL.Path)
	realIP:=fmt.Sprintf("RemoteAddr:[%s}, X-Forwarded-For:[%v] , X-Real-Ip:[%v]\n",r.RemoteAddr,
		r.Header.Get("X-Forwarded-For"),r.Header.Get("X-Real-Ip"))
	// 通过io包的WriteString中函数，往ResponseWriter中写会响应。
	io.WriteString(w, path)
	io.WriteString(w,realIP)
}

func (rs *RealServer) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("server error")
	w.WriteHeader(500)
	io.WriteString(w, path)
}

func main() {
	s1:=RealServer{
		Addr: "127.0.0.1:8081",
	}
	s1.RUN()
	s2:=RealServer{
		Addr: "127.0.0.1:8082",
	}
	s2.RUN()

	// 阻塞main协程
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan,syscall.SIGINT,syscall.SIGTERM)
	<-signalChan
}
