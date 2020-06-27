package main

import (
	"fmt"
	"net"
)

func main() {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 9091,
	})
	defer udpConn.Close()
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
	for {
		bytes := make([]byte, 1024)
		num, addr, err := udpConn.ReadFromUDP(bytes)
		if err != nil {
			fmt.Printf("error : %v", err)
			return
		}
		fmt.Printf("From remote address:[%v] read:[%v]byte content:[%v]", addr, num, string(bytes))
		// 回复消息
		go response(udpConn, addr)
	}
}

// 回复数据
func response(conn *net.UDPConn, addr *net.UDPAddr) {
	res := []byte("recieved your mag")
	_, err := conn.WriteToUDP(res, addr)
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
}
