package main

import (
	"fmt"
	"net"
)

func main() {

	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9091,
	})
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}

	_, err = udpConn.Write([]byte("i am udp client"))
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
	bytes:=make([]byte,1024)
	num, addr, err := udpConn.ReadFromUDP(bytes)
	if err != nil {
		fmt.Printf("Fail to read from udp error: [%v]", err)
		return
	}
	fmt.Printf("Recieve from udp address:[%v], bytes:[%v], content:[%v]",addr,num,string(bytes))

}
